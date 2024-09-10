package db

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/crypto"
	"luna-backend/sources"
	"luna-backend/sources/caldav"
	"luna-backend/types"
	"strings"

	"github.com/google/uuid"
)

type sourceEntry struct {
	Id       types.ID
	Name     string
	Type     string
	Settings string
}

func (db *Database) initializeSourcesTable() error {
	var err error
	// Sources table:
	// id user name type settings auth
	_, err = db.connection.Exec(`
		CREATE TABLE IF NOT EXISTS sources (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id),
			name VARCHAR(255) NOT NULL,
			type SOURCE_TYPE_ENUM NOT NULL,
			settings JSONB NOT NULL,
			auth_type BYTEA NOT NULL,
			auth BYTEA NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create sources table: %v", err)
	}

	return nil
}

func (db *Database) parseSource(rows types.PgxScanner) (sources.Source, error) {
	var err error
	var authType string
	var authBytes string
	sourceEntry := sourceEntry{}
	err = rows.Scan(&sourceEntry.Id, &sourceEntry.Name, &sourceEntry.Type, &sourceEntry.Settings, &authType, &authBytes)
	if err != nil {
		return nil, fmt.Errorf("could not scan source row: %v", err)
	}

	var authMethod auth.AuthMethod
	switch authType {
	case "none":
		authMethod = auth.NewNoAuth()
	case "basic":
		basicAuth := &auth.BasicAuth{}
		err = json.Unmarshal([]byte(authBytes), basicAuth)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal basic auth: %v", err)
		}
		authMethod = basicAuth
	case "bearer":
		bearerAuth := &auth.BearerAuth{}
		err = json.Unmarshal([]byte(authBytes), bearerAuth)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal bearer auth: %v", err)
		}
		authMethod = bearerAuth
	default:
		return nil, fmt.Errorf("unknown auth type: %v", authType)
	}

	switch sourceEntry.Type {
	case "caldav":
		settings := &caldav.CaldavSettings{}
		err = json.Unmarshal([]byte(sourceEntry.Settings), settings)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal caldav settings: %v", err)
		}
		caldavSource := caldav.PackCaldavSource(
			sourceEntry.Id,
			sourceEntry.Name,
			settings,
			authMethod,
		)
		return caldavSource, nil
	case "ical":
		fallthrough
	default:
		return nil, fmt.Errorf("unknown source type: %v", sourceEntry.Type)
	}
}

func getUserEncryptionKey(userId types.ID) (string, error) {
	masterKey, err := crypto.GetSymmetricKey("database")
	if err != nil {
		return "", fmt.Errorf("could not get master key: %v", err)
	}
	userKey, err := crypto.DeriveKey(masterKey, userId.Bytes())
	if err != nil {
		return "", fmt.Errorf("could not derive user key: %v", err)
	}
	encodedKey := base64.StdEncoding.EncodeToString(userKey)
	return encodedKey, nil
}

func getUserDecryptionKey(userId types.ID) (string, error) {
	return getUserEncryptionKey(userId)
}

func (db *Database) GetSource(userId types.ID, sourceId types.ID) (sources.Source, error) {
	decryptionKey, err := getUserDecryptionKey(userId)
	if err != nil {
		return nil, fmt.Errorf("could not get user decryption key: %v", err)
	}

	row := db.connection.QueryRow(`
		SELECT id, name, type, settings, PGP_SYM_DECRYPT(auth_type, $3), PGP_SYM_DECRYPT(auth, $3)
		FROM sources
		WHERE user_id = $1 AND id = $2;
	`, userId.UUID(), sourceId.UUID(), decryptionKey)

	source, err := db.parseSource(row)
	if err != nil {
		return nil, fmt.Errorf("could not parse source: %v", err)
	}

	return source, nil
}

func (db *Database) GetSources(userId types.ID) ([]sources.Source, error) {
	var err error

	decryptionKey, err := getUserDecryptionKey(userId)
	if err != nil {
		return nil, fmt.Errorf("could not get user decryption key: %v", err)
	}

	rows, err := db.connection.Query(`
		SELECT id, name, type, settings, PGP_SYM_DECRYPT(auth_type, $2), PGP_SYM_DECRYPT(auth, $2)
		FROM sources
		WHERE user_id = $1;
	`, userId.UUID(), decryptionKey)
	if err != nil {
		return nil, fmt.Errorf("could not get sources: %v", err)
	}
	defer rows.Close()

	sources := []sources.Source{}
	for rows.Next() {
		source, err := db.parseSource(rows)
		if err != nil {
			return nil, fmt.Errorf("could not parse source: %v", err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (db *Database) InsertSource(userId types.ID, source sources.Source) (types.ID, error) {
	encryptionKey, err := getUserEncryptionKey(userId)
	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not get user encryption key: %v", err)
	}

	query := `
		INSERT INTO sources (user_id, name, type, settings, auth_type, auth)
		VALUES ($1, $2, $3, $4, PGP_SYM_ENCRYPT($5, $7), PGP_SYM_ENCRYPT($6, $7))
		RETURNING id;
	`
	marshalledAuth, err := source.GetAuth().String()
	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not marshal auth: %v", err)
	}
	args := []any{userId.UUID(), source.GetName(), source.GetType(), source.GetSettings(), source.GetAuth().GetType(), marshalledAuth, encryptionKey}

	var id uuid.UUID
	err = db.connection.QueryRow(query, args...).Scan(&id)

	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not insert source: %v", err)
	}

	return types.IdFromUuid(id), nil
}

func (db *Database) UpdateSource(userId types.ID, sourceId types.ID, newName string, newAuth auth.AuthMethod, newSourceType string, newSourceSettings sources.SourceSettings) error {
	encryptionKey, err := getUserEncryptionKey(userId)
	if err != nil {
		return fmt.Errorf("could not get user encryption key: %v", err)
	}

	changes := []string{}
	args := []any{}

	if newName != "" {
		changes = append(changes, fmt.Sprintf("name = $%d", len(changes)+1))
		args = append(args, newName)
	}
	if newSourceType != "" {
		changes = append(changes, fmt.Sprintf("type = $%d", len(changes)+1), fmt.Sprintf("settings = $%d", len(changes)+2))
		args = append(args, newSourceType, newSourceSettings)
	}
	if newAuth != nil {
		changes = append(
			changes,
			fmt.Sprintf("auth_type = PGP_SYM_ENCRYPT($%d, $%d)", len(changes)+1, len(changes)+3),
			fmt.Sprintf("auth = PGP_SYM_ENCRYPT($%d, $%d)", len(changes)+2, len(changes)+3),
		)
		marshalledAuth, err := newAuth.String()
		if err != nil {
			return fmt.Errorf("could not marshal auth: %v", err)
		}
		args = append(args, newAuth.GetType(), marshalledAuth)
	}

	if len(changes) == 0 {
		return fmt.Errorf("no changes to update")
	}

	query := fmt.Sprintf(`
		UPDATE sources
		SET %s
		WHERE user_id = $%d AND id = $%d;
	`, strings.Join(changes, ", "), len(changes)+2, len(changes)+3)
	args = append(args, encryptionKey, userId.UUID(), sourceId.UUID())

	_, err = db.connection.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("could not update source: %v", err)
	}

	return nil
}

func (db *Database) DeleteSource(userId types.ID, sourceId types.ID) error {
	tag, err := db.connection.Exec(`
		DELETE FROM sources
		WHERE user_id = $1 AND id = $2;
	`, userId.UUID(), sourceId)
	if err != nil {
		return fmt.Errorf("could not delete source: %v", err)
	}
	if tag.RowsAffected() == 0 {
		// TODO: consider not returning an error here
		// TODO: this is of essence on lossy networks
		// TODO: if the first delete confirmation fails and the user retries,
		// TODO: we can simply confirm that the source no longer exists
		return fmt.Errorf("could not delete source: source not found")
	}

	return nil
}

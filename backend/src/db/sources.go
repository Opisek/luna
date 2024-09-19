package db

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/common"
	"luna-backend/crypto"
	"luna-backend/interface/primitives"
	"luna-backend/interface/protocols/caldav"
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

func (tx *Transaction) initializeSourcesTable() error {
	var err error
	// Sources table:
	// id user name type settings auth
	_, err = tx.conn.Exec(
		context.TODO(),
		`
		CREATE TABLE IF NOT EXISTS sources (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id),
			name VARCHAR(255) NOT NULL,
			type SOURCE_TYPE_ENUM NOT NULL,
			settings JSONB NOT NULL,
			auth_type BYTEA NOT NULL,
			auth BYTEA NOT NULL,
			UNIQUE (user_id, name)
		);
		`,
	)
	if err != nil {
		return fmt.Errorf("could not create sources table: %v", err)
	}

	return nil
}

func parseSource(rows types.PgxScanner) (primitives.Source, error) {
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
	case types.AuthNone:
		authMethod = auth.NewNoAuth()
	case types.AuthBasic:
		basicAuth := &auth.BasicAuth{}
		err = json.Unmarshal([]byte(authBytes), basicAuth)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal basic auth: %v", err)
		}
		authMethod = basicAuth
	case types.AuthBearer:
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
	case types.SourceCaldav:
		settings := &caldav.CaldavSourceSettings{}
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
	case types.SourceIcal:
		fallthrough
	default:
		return nil, fmt.Errorf("unknown source type: %v", sourceEntry.Type)
	}
}

func getUserEncryptionKey(commonConfig *common.CommonConfig, userId types.ID) (string, error) {
	masterKey, err := crypto.GetSymmetricKey(commonConfig, "database")
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

func getUserDecryptionKey(commonConfig *common.CommonConfig, userId types.ID) (string, error) {
	return getUserEncryptionKey(commonConfig, userId)
}

func (tx *Transaction) GetSource(userId types.ID, sourceId types.ID) (primitives.Source, error) {
	decryptionKey, err := getUserDecryptionKey(tx.db.commonConfig, userId)
	if err != nil {
		return nil, fmt.Errorf("could not get user decryption key: %v", err)
	}

	row := tx.conn.QueryRow(
		context.TODO(),
		`
		SELECT id, name, type, settings, PGP_SYM_DECRYPT(auth_type, $3), PGP_SYM_DECRYPT(auth, $3)
		FROM sources
		WHERE user_id = $1 AND id = $2;
		`,
		userId.UUID(),
		sourceId.UUID(),
		decryptionKey,
	)

	source, err := parseSource(row)
	if err != nil {
		return nil, fmt.Errorf("could not parse source: %v", err)
	}

	return source, nil
}

func (tx *Transaction) GetSources(userId types.ID) ([]primitives.Source, error) {
	var err error

	decryptionKey, err := getUserDecryptionKey(tx.db.commonConfig, userId)
	if err != nil {
		return nil, fmt.Errorf("could not get user decryption key: %v", err)
	}

	rows, err := tx.conn.Query(
		context.TODO(),
		`
		SELECT id, name, type, settings, PGP_SYM_DECRYPT(auth_type, $2), PGP_SYM_DECRYPT(auth, $2)
		FROM sources
		WHERE user_id = $1;
		`,
		userId.UUID(),
		decryptionKey,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get sources: %v", err)
	}
	defer rows.Close()

	sources := []primitives.Source{}
	for rows.Next() {
		source, err := parseSource(rows)
		if err != nil {
			return nil, fmt.Errorf("could not parse source: %v", err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (tx *Transaction) InsertSource(userId types.ID, source primitives.Source) (types.ID, error) {
	encryptionKey, err := getUserEncryptionKey(tx.db.commonConfig, userId)
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
	err = tx.conn.QueryRow(context.TODO(), query, args...).Scan(&id)

	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not insert source: %v", err)
	}

	return types.IdFromUuid(id), nil
}

func (tx *Transaction) UpdateSource(userId types.ID, sourceId types.ID, newName string, newAuth auth.AuthMethod, newSourceType string, newSourceSettings primitives.SourceSettings) error {
	encryptionKey, err := getUserEncryptionKey(tx.db.commonConfig, userId)
	if err != nil {
		return fmt.Errorf("could not get user encryption key: %v", err)
	}

	changes := []string{}
	args := []any{}

	if newName != "" {
		changes = append(changes, fmt.Sprintf("name = $%d", len(changes)+1))
		args = append(args, newName)
	}
	if newSourceType != "" && newSourceSettings != nil {
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
		args = append(args, newAuth.GetType(), marshalledAuth, encryptionKey)
	}

	if len(changes) == 0 {
		return fmt.Errorf("no changes to update")
	}

	query := fmt.Sprintf(`
		UPDATE sources
		SET %s
		WHERE user_id = $%d AND id = $%d;
	`, strings.Join(changes, ", "), len(args)+1, len(args)+2)
	args = append(args, userId.UUID(), sourceId.UUID())

	_, err = tx.conn.Exec(context.TODO(), query, args...)

	if err != nil {
		return fmt.Errorf("could not update source: %v", err)
	}

	return nil
}

func (tx *Transaction) DeleteSource(userId types.ID, sourceId types.ID) (bool, error) {
	tag, err := tx.conn.Exec(
		context.TODO(),
		`
		DELETE FROM sources
		WHERE user_id = $1 AND id = $2;
		`,
		userId.UUID(),
		sourceId,
	)
	if err != nil {
		return false, fmt.Errorf("could not delete source: %v", err)
	}
	return tag.RowsAffected() != 0, nil
}

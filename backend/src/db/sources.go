package db

import (
	"encoding/json"
	"fmt"
	"luna-backend/auth"
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
	// TODO: add auth_type and auth
	_, err = db.connection.Exec(`
		CREATE TABLE IF NOT EXISTS sources (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id),
			name VARCHAR(255) NOT NULL,
			type SOURCE_TYPE_ENUM NOT NULL,
			settings JSONB NOT NULL,
			auth_type AUTH_TYPE_ENUM NOT NULL,
			auth JSONB NOT NULL
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
	var authBytes []byte
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
		err = json.Unmarshal(authBytes, basicAuth)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal basic auth: %v", err)
		}
		authMethod = basicAuth
	case "bearer":
		bearerAuth := &auth.BearerAuth{}
		err = json.Unmarshal(authBytes, bearerAuth)
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

func (db *Database) GetSource(userId types.ID, sourceId types.ID) (sources.Source, error) {
	db.logger.Debugf("Getting source %v for user %v", sourceId, userId)

	row := db.connection.QueryRow(`
		SELECT id, name, type, settings, auth_type, auth
		FROM sources
		WHERE user_id = $1 AND id = $2;
	`, userId.UUID(), sourceId.UUID())

	db.logger.Debugf("Got source %v for user %v", sourceId, userId)

	source, err := db.parseSource(row)
	if err != nil {
		return nil, fmt.Errorf("could not parse source: %v", err)
	}

	db.logger.Debugf("Parsed source %v for user %v", sourceId, userId)

	return source, nil
}

func (db *Database) GetSources(userId types.ID) ([]sources.Source, error) {
	var err error

	rows, err := db.connection.Query(`
		SELECT id, name, type, settings, auth_type, auth
		FROM sources
		WHERE user_id = $1;
	`, userId.UUID())
	if err != nil {
		return nil, err
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
	query := `
		INSERT INTO sources (user_id, name, type, settings, auth_type, auth)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	// TODO: encrypt auth before saving in the database
	args := []any{userId.UUID(), source.GetName(), source.GetType(), source.GetSettings(), source.GetAuth().GetType(), source.GetAuth()}

	var id uuid.UUID
	err := db.connection.QueryRow(query, args...).Scan(&id)

	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not insert source: %v", err)
	}

	return types.IdFromUuid(id), nil
}

func (db *Database) UpdateSource(userId types.ID, sourceId types.ID, newName string, newAuth auth.AuthMethod, newSourceType string, newSourceSettings sources.SourceSettings) error {
	changes := []string{}
	args := []any{}

	if newName != "" {
		changes = append(changes, fmt.Sprintf("name = $%d", len(changes)+1))
		args = append(args, newName)
	}
	if newAuth != nil {
		changes = append(changes, fmt.Sprintf("auth = $%d", len(changes)+1))
		args = append(args, newAuth)
	}
	if newSourceType != "" {
		changes = append(changes, fmt.Sprintf("type = $%d", len(changes)+1), fmt.Sprintf("settings = $%d", len(changes)+2))
		args = append(args, newSourceType, newSourceSettings)
	}

	query := fmt.Sprintf(`
		UPDATE sources
		SET %s
		WHERE user_id = $%d AND id = $%d;
	`, strings.Join(changes, ", "), len(changes)+1, len(changes)+2)
	args = append(args, userId.UUID(), sourceId.UUID())

	_, err := db.connection.Exec(query, args...)

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

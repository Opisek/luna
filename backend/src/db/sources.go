package db

import (
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/sources"
	"luna-backend/sources/caldav"
	"luna-backend/types"
	"strings"
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
			type source_type NOT NULL,
			settings JSONB NOT NULL
			auth_type auth_type NOT NULL,
			auth JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create sources table: %v", err)
	}

	return nil
}

func (db *Database) GetSources(userId types.ID) ([]sources.Source, error) {
	var err error

	// TODO: also get auth_type and auth?
	rows, err := db.connection.Query(`
		SELECT id, name, type, settings
		FROM sources
		WHERE user_id = $1;
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := []sources.Source{}
	for rows.Next() {
		sourceEntry := sourceEntry{}
		err = rows.Scan(&sourceEntry.Id, &sourceEntry.Name, &sourceEntry.Type, &sourceEntry.Settings)
		if err != nil {
			return nil, fmt.Errorf("could not scan source row: %v", err)
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
				nil,
			)
			sources = append(sources, caldavSource)
		case "ical":
			fallthrough
		default:
			return nil, fmt.Errorf("unknown source type: %v", sourceEntry.Type)
		}
	}

	return sources, nil
}

func (db *Database) InsertSource(userId types.ID, source sources.Source) (types.ID, error) {
	query := `
		INSERT INTO sources (user_id, name, type, settings, auth_type, auth)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	args := []any{userId, source.GetName(), source.GetType(), source.GetSettings(), source.GetAuth().GetType(), source.GetAuth()}

	var id types.ID
	err := db.connection.QueryRow(query, args...).Scan(&id)

	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not insert source: %v", err)
	}

	return id, nil
}

func (db *Database) UpdateSource(userId types.ID, sourceId types.ID, newName string, newAuth auth.AuthMethod, newSourceType string, newSourceSettings []byte) error {
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
	args = append(args, userId, sourceId)

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
	`, userId, sourceId)
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

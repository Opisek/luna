package db

import (
	"encoding/json"
	"fmt"
	"luna-backend/sources"
	"luna-backend/sources/caldav"

	"github.com/google/uuid"
)

type sourceEntry struct {
	Id       sources.SourceId
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
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create sources table: %v", err)
	}

	return nil
}

func (db *Database) GetSources(userId uuid.UUID) ([]sources.Source, error) {
	var err error

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

func (db *Database) InsertSource(userId uuid.UUID, source sources.Source) error {
	query := `
		INSERT INTO sources (user_id, name, type, settings)
		VALUES ($1, $2, $3, $4);
	`
	args := []any{userId, source.GetName(), source.GetType(), source.GetSettings()}

	_, err := db.connection.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("could not insert source: %v", err)
	}

	return nil
}

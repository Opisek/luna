package db

import (
	"fmt"
	"luna-backend/interface/primitives"
	"luna-backend/types"

	"github.com/jackc/pgx"
)

type calendarEntry struct {
	Id       types.ID
	Source   types.ID
	Color    *types.Color
	Settings primitives.CalendarSettings
}

func (db *Database) initializeCalendarsTable() error {
	var err error
	// Calendars table:
	// id source color settings
	_, err = db.connection.Exec(`
		CREATE TABLE IF NOT EXISTS calendars (
			id UUID PRIMARY KEY,
			source UUID REFERENCES sources(id),
			color BYTEA,
			settings JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create calendars table: %v", err)
	}

	return nil
}

func (db *Database) insertCalendars(cals []primitives.Calendar) error {
	rows := [][]any{}

	for _, cal := range cals {
		row := []any{
			cal.GetId(),
			cal.GetSource().GetId(),
			cal.GetColor().Bytes(),
			cal.GetSettings().Bytes(),
		}

		rows = append(rows, row)
	}

	// TODO: to avoid conflicts with existing keys, we want to do something similar to this:
	// TODO: https://github.com/jackc/pgx/issues/992
	// TODO: this might require transactions to be set up first
	_, err := db.connection.CopyFrom(
		pgx.Identifier{"calendars"},
		[]string{"id", "source", "color", "settings"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return fmt.Errorf("could not insert calendars into database: %v", err)
	}

	return nil
}

func (db *Database) GetCalendars(source primitives.Source) ([]primitives.Calendar, error) {
	cals, err := source.GetCalendars()
	if err != nil {
		return nil, fmt.Errorf("could not get calendars from source %v: %v", source.GetId().String(), err)
	}

	err = db.insertCalendars(cals)
	if err != nil {
		return nil, fmt.Errorf("could not cache calendars: %v", err)
	}

	return cals, nil
}

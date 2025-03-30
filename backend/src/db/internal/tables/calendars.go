package tables

import (
	"fmt"
)

func (q *Tables) InitializeCalendarsTable() error {
	var err error
	// Calendars table:
	// id source settings
	_, err = q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS calendars (
			id UUID PRIMARY KEY,
			source UUID REFERENCES sources(id) ON DELETE CASCADE,
			settings JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create calendars table: %v", err)
	}

	return nil
}

func (q *Tables) InitializeCalendarOverridesTable() error {
	var err error
	// Calendars table:
	// id title description color
	_, err = q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS calendar_overrides (
			calendarid UUID UNIQUE REFERENCES calendars(id) ON DELETE CASCADE,
			title TEXT,
			description TEXT,
			color BYTEA
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create calendar overrides table: %v", err)
	}

	return nil
}

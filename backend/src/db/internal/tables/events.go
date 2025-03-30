package tables

import (
	"fmt"
)

func (q *Tables) InitializeEventsTable() error {
	var err error
	// Events table:
	// id calendar settings
	_, err = q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY,
			calendar UUID REFERENCES calendars(id) ON DELETE CASCADE,
			settings JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create events table: %v", err)
	}

	return nil
}

func (q *Tables) InitializeEventOverridesTable() error {
	var err error
	// Calendars table:
	// id title description color
	_, err = q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS event_overrides (
			eventid UUID UNIQUE REFERENCES events(id) ON DELETE CASCADE,
			title TEXT,
			description TEXT,
			color BYTEA
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create event overrides table: %v", err)
	}

	return nil
}

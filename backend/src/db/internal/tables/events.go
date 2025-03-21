package tables

import (
	"fmt"
)

func (q *Tables) InitializeEventsTable() error {
	var err error
	// Events table:
	// id calendar color settings
	_, err = q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY,
			calendar UUID REFERENCES calendars(id) ON DELETE CASCADE,
			color BYTEA,
			settings JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create events table: %v", err)
	}

	return nil
}

package tables

import (
	"context"
	"fmt"
)

func (q *Tables) InitializeCalendarsTable() error {
	var err error
	// Calendars table:
	// id source color settings
	_, err = q.Tx.Exec(
		context.TODO(),
		`
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

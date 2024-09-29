package tables

import (
	"context"
	"fmt"
	"luna-backend/types"
)

type CalendarEntry struct {
	Id       types.ID `db:"id" encrypted:"false"`
	Source   types.ID `db:"source" encrypted:"false"`
	Color    []byte   `db:"color" encrypted:"false"`
	Settings []byte   `db:"settings" encrypted:"false"`
}

func (q *Tables) InitializeCalendarsTable() error {
	var err error
	// Calendars table:
	// id source color settings
	_, err = q.Tx.Exec(
		context.TODO(),
		`
		CREATE TABLE IF NOT EXISTS calendars (
			id UUID PRIMARY KEY,
			source UUID REFERENCES sources(id) ON DELETE CASCADE,
			color BYTEA,
			settings JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create calendars table: %v", err)
	}

	return nil
}

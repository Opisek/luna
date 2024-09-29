package tables

import (
	"context"
	"fmt"
	"luna-backend/types"
)

type EventEntry struct {
	Id       types.ID `db:"id" encrypted:"false"`
	Calendar types.ID `db:"calendar" encrypted:"false"`
	Color    []byte   `db:"color" encrypted:"false"`
	Settings []byte   `db:"settings" encrypted:"false"`
}

func (q *Tables) InitializeEventsTable() error {
	var err error
	// Events table:
	// id calendar color settings
	_, err = q.Tx.Exec(
		context.TODO(),
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

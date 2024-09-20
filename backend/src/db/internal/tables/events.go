package tables

import (
	"context"
	"fmt"
	"luna-backend/interface/primitives"
	"luna-backend/types"
)

type EventEntry struct {
	Id       types.ID
	Calendar primitives.Calendar
	Color    *types.Color
	Settings primitives.EventSettings
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
			calendar UUID REFERENCES calendars(id),
			color BYTEA,
			settings JSONB NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create events table: %v", err)
	}

	return nil
}

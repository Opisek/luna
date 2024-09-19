package db

import (
	"context"
	"encoding/json"
	"fmt"
	"luna-backend/interface/primitives"
	"luna-backend/interface/protocols/caldav"
	"luna-backend/types"
	"time"
)

type eventEntry struct {
	Id       types.ID
	Calendar primitives.Calendar
	Color    *types.Color
	Settings primitives.EventSettings
}

func (tx *Transaction) initializeEventsTable() error {
	var err error
	// Events table:
	// id calendar color settings
	_, err = tx.conn.Exec(
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

func (tx *Transaction) insertEvents(cals []primitives.Event) error {
	rows := [][]any{}

	for _, cal := range cals {
		color := cal.GetColor()
		var colBytes []byte
		if color.IsEmpty() {
			colBytes = nil
		} else {
			colBytes = color.Bytes()
		}

		row := []any{
			cal.GetId(),
			cal.GetCalendar().GetId(),
			colBytes,
			cal.GetSettings().Bytes(),
		}

		rows = append(rows, row)
	}

	err := tx.CopyAndUpdate(
		context.TODO(),
		"events",
		[]string{"id", "calendar", "color", "settings"},
		[]string{"color", "settings"},
		rows,
	)

	if err != nil {
		return fmt.Errorf("could not insert events into database: %v", err)
	}

	return nil
}

func parseEventSettings(sourceType string, settings []byte) (primitives.CalendarSettings, error) {
	switch sourceType {
	case types.SourceCaldav:
		parsedSettings := &caldav.CaldavEventSettings{}
		err := json.Unmarshal(settings, parsedSettings)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal caldav settings: %v", err)
		}
		return parsedSettings, nil
	case types.SourceIcal:
		fallthrough
	default:
		return nil, fmt.Errorf("unknown source type %v", sourceType)
	}
}

func (tx *Transaction) getEvents(calendar primitives.Calendar) ([]*eventEntry, error) {
	rows, err := tx.conn.Query(
		context.TODO(),
		`
		SELECT id, color, settings
		FROM events
		WHERE calendar = $1;
		`,
		calendar.GetId(),
	)

	if err != nil {
		return nil, fmt.Errorf("could not get events from database: %v", err)
	}

	defer rows.Close()

	cals := []*eventEntry{}
	for rows.Next() {
		var id types.ID
		var color []byte
		var settings []byte

		err := rows.Scan(&id, &color, &settings)
		if err != nil {
			return nil, fmt.Errorf("could not scan event row: %v", err)
		}

		parsedSettings, err := parseEventSettings(calendar.GetSource().GetType(), settings)
		if err != nil {
			return nil, fmt.Errorf("could not parse event settings: %v", err)
		}

		cals = append(cals, &eventEntry{
			Id:       id,
			Calendar: calendar,
			Color:    types.ColorFromBytes(color),
			Settings: parsedSettings,
		})
	}

	return cals, nil
}

func (tx *Transaction) GetEvents(calendar primitives.Calendar, start time.Time, end time.Time) ([]primitives.Event, error) {
	events, err := calendar.GetEvents(start, end)
	if err != nil {
		return nil, fmt.Errorf("could not get events from calendar %v: %v", calendar.GetId().String(), err)
	}

	eventMap := map[types.ID]primitives.Event{}
	for _, event := range events {
		eventMap[event.GetId()] = event
	}

	dbEvents, err := tx.getEvents(calendar)
	if err != nil {
		return nil, fmt.Errorf("could not get cached events: %v", err)
	}

	for _, dbEvent := range dbEvents {
		if event, ok := eventMap[dbEvent.Id]; ok {
			if event.GetColor() == nil {
				event.SetColor(dbEvent.Color)
			}
		}
	}

	err = tx.insertEvents(events)
	if err != nil {
		return nil, fmt.Errorf("could not cache events: %v", err)
	}

	return events, nil
}

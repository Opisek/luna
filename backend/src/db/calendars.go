package db

import (
	"context"
	"encoding/json"
	"fmt"
	"luna-backend/interface/primitives"
	"luna-backend/interface/protocols/caldav"
	"luna-backend/types"
)

type calendarEntry struct {
	Id       types.ID
	Source   primitives.Source
	Color    *types.Color
	Settings primitives.CalendarSettings
}

func (tx *Transaction) initializeCalendarsTable() error {
	var err error
	// Calendars table:
	// id source color settings
	_, err = tx.conn.Exec(
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

func (tx *Transaction) insertCalendars(cals []primitives.Calendar) error {
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

	err := tx.CopyAndUpdate(
		context.TODO(),
		"calendars",
		[]string{"id", "source", "color", "settings"},
		[]string{"color", "settings"},
		rows,
	)

	if err != nil {
		return fmt.Errorf("could not insert calendars into database: %v", err)
	}

	return nil
}

func parseCalendarSettings(sourceType string, settings []byte) (primitives.CalendarSettings, error) {
	switch sourceType {
	case types.SourceCaldav:
		parsedSettings := &caldav.CaldavCalendarSettings{}
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

func (tx *Transaction) getCalendars(source primitives.Source) ([]*calendarEntry, error) {
	rows, err := tx.conn.Query(
		context.TODO(),
		`
		SELECT id, color, settings
		FROM calendars
		WHERE source = $1;
	`, source.GetId())

	if err != nil {
		return nil, fmt.Errorf("could not get calendars from database: %v", err)
	}

	defer rows.Close()

	cals := []*calendarEntry{}
	for rows.Next() {
		var id types.ID
		var color []byte
		var settings []byte

		err := rows.Scan(&id, &color, &settings)
		if err != nil {
			return nil, fmt.Errorf("could not scan calendar row: %v", err)
		}

		parsedSettings, err := parseCalendarSettings(source.GetType(), settings)
		if err != nil {
			return nil, fmt.Errorf("could not parse calendar settings: %v", err)
		}

		cals = append(cals, &calendarEntry{
			Id:       id,
			Source:   source,
			Color:    types.ColorFromBytes(color),
			Settings: parsedSettings,
		})
	}

	return cals, nil
}

func (tx *Transaction) GetCalendars(source primitives.Source) ([]primitives.Calendar, error) {
	cals, err := source.GetCalendars()
	if err != nil {
		return nil, fmt.Errorf("could not get calendars from source %v: %v", source.GetId().String(), err)
	}

	calMap := map[types.ID]primitives.Calendar{}
	for _, cal := range cals {
		calMap[cal.GetId()] = cal
	}

	dbCals, err := tx.getCalendars(source)
	if err != nil {
		return nil, fmt.Errorf("could not get cached calendars: %v", err)
	}

	for _, dbCal := range dbCals {
		if cal, ok := calMap[dbCal.Id]; ok {
			if cal.GetColor() == nil {
				cal.SetColor(dbCal.Color)
				// TODO: if dbCal.Color == nil, either return some default color, or generate a deterministic random one (e.g. calendar id hash -> hue)
			}
		}
	}

	err = tx.insertCalendars(cals)
	if err != nil {
		return nil, fmt.Errorf("could not cache calendars: %v", err)
	}

	return cals, nil
}

func (tx *Transaction) UpdateCalendar(cal primitives.Calendar) error {
	_, err := tx.conn.Exec(
		context.TODO(),
		`
		UPDATE calendars
		SET color = $1, settings = $2
		WHERE id = $3;`,
		cal.GetColor().Bytes(),
		cal.GetSettings().Bytes(),
		cal.GetId(),
	)

	if err != nil {
		return fmt.Errorf("could not update calendar %v: %v", cal.GetId().String(), err)
	}

	return nil
}

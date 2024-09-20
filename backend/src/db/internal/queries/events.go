package queries

import (
	"context"
	"fmt"
	"luna-backend/db/internal/parsing"
	"luna-backend/db/internal/tables"
	"luna-backend/db/internal/util"
	"luna-backend/interface/primitives"
	"luna-backend/types"
)

func (q *Queries) insertEvents(cals []primitives.Event) error {
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

	err := util.CopyAndUpdate(
		q.Tx,
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

func (q *Queries) getEventEntries(calendars []primitives.Calendar) ([]*tables.EventEntry, error) {
	query := fmt.Sprintf(
		`
		SELECT id, calendar, color, settings
		FROM events
		WHERE calendar IN (
			%s
		);
		`,
		util.GenerateArgList(1, len(calendars)),
	)

	rows, err := q.Tx.Query(
		context.TODO(),
		query,
		util.JoinIds(calendars, func(s primitives.Calendar) types.ID { return s.GetId() })...,
	)

	if err != nil {
		return nil, fmt.Errorf("could not get calendars from database: %v", err)
	}

	defer rows.Close()

	calMap := map[types.ID]primitives.Calendar{}
	for _, cal := range calendars {
		calMap[cal.GetId()] = cal
	}

	events := []*tables.EventEntry{}
	for rows.Next() {
		var id types.ID
		var source types.ID
		var color []byte
		var settings []byte

		err := rows.Scan(&id, &source, &color, &settings)
		if err != nil {
			return nil, fmt.Errorf("could not scan calendar row: %v", err)
		}

		eventEntry, err := parsing.ParseEventEntry(calMap[source], id, color, settings)
		if err != nil {
			return nil, fmt.Errorf("could not parse calendar: %v", err)
		}

		events = append(events, eventEntry)
	}

	return events, nil
}

func (q *Queries) ReconcileEvents(cals []primitives.Calendar, events []primitives.Event) ([]primitives.Event, error) {
	eventMap := map[types.ID]primitives.Event{}
	for _, event := range events {
		eventMap[event.GetId()] = event
	}

	dbEvents, err := q.getEventEntries(cals)
	if err != nil {
		return nil, fmt.Errorf("could not get cached events: %v", err)
	}

	for _, dbEvent := range dbEvents {
		if event, ok := eventMap[dbEvent.Id]; ok {
			if event.GetColor() == nil {
				event.SetColor(dbEvent.Color)
				// TODO: if dbCal.Color == nil, either return some default color, or generate a deterministic random one (e.g. calendar id hash -> hue)
			}
		}
	}

	err = q.insertEvents(events)
	if err != nil {
		return nil, fmt.Errorf("could not cache events: %v", err)
	}

	return events, nil
}

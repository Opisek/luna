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

func (q *Queries) getEventEntries(calendars []primitives.Calendar, events []primitives.Event) ([]*tables.EventEntry, error) {
	query := fmt.Sprintf(
		`
		SELECT id, calendar, color, settings
		FROM events
		WHERE id IN (
			%s
		);
		`,
		util.GenerateArgList(1, len(events)),
	)

	rows, err := q.Tx.Query(
		context.TODO(),
		query,
		util.JoinIds(events, func(e primitives.Event) types.ID { return e.GetId() })...,
	)

	if err != nil {
		return nil, fmt.Errorf("could not get calendars from database: %v", err)
	}

	defer rows.Close()

	calMap := map[types.ID]primitives.Calendar{}
	for _, cal := range calendars {
		calMap[cal.GetId()] = cal
	}

	eventEntries := []*tables.EventEntry{}
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

		eventEntries = append(eventEntries, eventEntry)
	}

	return eventEntries, nil
}

func (q *Queries) ReconcileEvents(cals []primitives.Calendar, events []primitives.Event) ([]primitives.Event, error) {
	eventMap := map[types.ID]primitives.Event{}
	for _, event := range events {
		eventMap[event.GetId()] = event
	}

	dbEvents, err := q.getEventEntries(cals, events)
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

func (q *Queries) GetEvent(userId types.ID, eventId types.ID) (primitives.Event, error) {
	var err error

	decryptionKey, err := util.GetUserDecryptionKey(q.CommonConfig, userId)
	if err != nil {
		return nil, fmt.Errorf("could not get user decryption key: %v", err)
	}

	var sourceId types.ID
	var sourceName string
	var sourceType string
	var sourceSettings []byte
	var authType string
	var authBytes []byte
	var calendarId types.ID
	var calendarColor []byte
	var calendarSettings []byte
	var eventColor []byte
	var eventSettings []byte

	err = q.Tx.QueryRow(
		context.TODO(),
		`
		SELECT sources.id, sources.name, sources.type, sources.settings, PGP_SYM_DECRYPT(sources.auth_type, $3), PGP_SYM_DECRYPT(sources.auth, $3), calendars.id, calendars.color, calendars.settings, events.color, events.settings
		FROM events
		JOIN calendars ON events.calendar = calendars.id
		JOIN sources ON calendars.source = sources.id
		WHERE events.id = $1
		AND sources.userid = $2;
		`,
		eventId.UUID(),
		userId.UUID(),
		decryptionKey,
	).Scan(&sourceId, &sourceName, &sourceType, &sourceSettings, &authType, &authBytes, &calendarId, &calendarColor, &calendarSettings, &eventColor, &eventSettings)
	if err != nil {
		return nil, fmt.Errorf("could not get event: %v", err)
	}

	source, err := parsing.ParseSource(sourceId, sourceName, sourceType, sourceSettings, authType, authBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse source: %v", err)
	}

	calendarEntry, err := parsing.ParseCalendarEntry(source, calendarId, calendarColor, calendarSettings)
	if err != nil {
		return nil, fmt.Errorf("could not parse calendar entry: %v", err)
	}

	calendar, err := source.GetCalendar(calendarEntry.Settings)
	if err != nil {
		return nil, fmt.Errorf("could not get calendar: %v", err)
	}
	if calendar.GetColor() == nil {
		calendar.SetColor(calendarEntry.Color)
	}

	eventEntry, err := parsing.ParseEventEntry(calendar, eventId, eventColor, eventSettings)
	if err != nil {
		return nil, fmt.Errorf("could not parse event entry: %v", err)
	}

	event, err := calendar.GetEvent(eventEntry.Settings)
	if err != nil {
		return nil, fmt.Errorf("could not get event: %v", err)
	}
	if event.GetColor() == nil {
		event.SetColor(eventEntry.Color)
	}

	return event, nil
}

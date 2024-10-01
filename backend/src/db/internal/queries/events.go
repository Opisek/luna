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

func (q *Queries) getEventEntries(events []primitives.Event) ([]*tables.EventEntry, error) {
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

	entries := []*tables.EventEntry{}
	for rows.Next() {
		entry := &tables.EventEntry{}

		err := rows.Scan(&entry.Id, &entry.Calendar, &entry.Color, &entry.Settings)
		if err != nil {
			return nil, fmt.Errorf("could not scan calendar row: %v", err)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (q *Queries) ReconcileEvents(events []primitives.Event) ([]primitives.Event, error) {
	if len(events) == 0 {
		return events, nil
	}

	eventMap := map[types.ID]primitives.Event{}
	for _, event := range events {
		eventMap[event.GetId()] = event
	}

	dbEvents, err := q.getEventEntries(events)
	if err != nil {
		return nil, fmt.Errorf("could not get cached events: %v", err)
	}

	for _, dbEvent := range dbEvents {
		if event, ok := eventMap[dbEvent.Id]; ok {
			if event.GetColor() == nil {
				event.SetColor(types.ColorFromBytes(dbEvent.Color))
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

	scanner := parsing.NewPgxScanner()
	scanner.ScheduleEvent()
	cols, params := scanner.Variables(3)

	query := fmt.Sprintf(
		`
		SELECT %s 
		FROM events
		JOIN calendars ON events.calendar = calendars.id
		JOIN sources ON calendars.source = sources.id
		WHERE events.id = $1
		AND sources.userid = $2;
		`,
		cols,
	)

	err = q.Tx.QueryRow(
		context.TODO(),
		query,
		eventId.UUID(),
		userId.UUID(),
		decryptionKey,
	).Scan(params...)
	if err != nil {
		return nil, fmt.Errorf("could not get event: %v", err)
	}

	return scanner.GetEvent()
}

func (q *Queries) InsertEvent(event primitives.Event) error {
	_, err := q.Tx.Exec(
		context.TODO(),
		`
		INSERT INTO events (id, calendar, color, settings)
		VALUES ($1, $2, $3, $4);
		`,
		event.GetId().UUID(),
		event.GetCalendar().GetId().UUID(),
		event.GetColor().Bytes(),
		event.GetSettings().Bytes(),
	)
	if err != nil {
		return fmt.Errorf("could not insert event into database: %v", err)
	}

	return nil
}

func (q *Queries) UpdateEvent(event primitives.Event) error {
	_, err := q.Tx.Exec(
		context.TODO(),
		`
		UPDATE events
		SET color = $2, settings = $3
		WHERE id = $1;
		`,
		event.GetId().UUID(),
		event.GetColor().Bytes(),
		event.GetSettings().Bytes(),
	)
	if err != nil {
		return fmt.Errorf("could not update event in database: %v", err)
	}

	return nil
}

func (q *Queries) DeleteEvent(userId types.ID, eventId types.ID) error {
	_, err := q.Tx.Exec(
		context.TODO(),
		`
		DELETE FROM events
		WHERE id = $1
		AND calendar IN (
			SELECT calendars.id
			FROM calendars
			JOIN sources ON calendars.source = sources.id
			WHERE sources.userid = $2
		);
		`,
		eventId.UUID(),
		userId.UUID(),
	)
	if err != nil {
		return fmt.Errorf("could not delete event from database: %v", err)
	}

	return nil
}

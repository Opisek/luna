package queries

import (
	"fmt"
	"luna-backend/db/internal/parsing"
	"luna-backend/db/internal/util"
	"luna-backend/errors"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) insertEvents(events []primitives.Event) *errors.ErrorTrace {
	rows := [][]any{}

	for _, event := range events {
		color := event.GetColor()
		var colBytes []byte
		if color.IsEmpty() {
			colBytes = nil
		} else {
			colBytes = color.Bytes()
		}

		row := []any{
			event.GetId(),
			event.GetCalendar().GetId(),
			colBytes,
			event.GetSettings().Bytes(),
		}

		rows = append(rows, row)
	}

	err := util.CopyAndUpdate(
		q.Tx,
		q.Context,
		"events",
		[]string{"id", "calendar", "color", "settings"},
		[]string{"color", "settings"},
		rows,
	)

	if err != nil {
		return err.
			Append(errors.LvlWordy, "Could not insert events")
	}

	return nil
}

func (q *Queries) getEventEntries(events []primitives.Event) ([]*types.EventDatabaseEntry, *errors.ErrorTrace) {
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
		q.Context,
		query,
		util.JoinIds(events, func(e primitives.Event) types.ID { return e.GetId() })...,
	)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get events from the database")
	}
	defer rows.Close()

	entries := []*types.EventDatabaseEntry{}
	for rows.Next() {
		entry := &types.EventDatabaseEntry{}

		err := rows.Scan(&entry.Id, &entry.Calendar, &entry.Color, &entry.Settings)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not scan event row")
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (q *Queries) ReconcileEvents(events []primitives.Event) ([]primitives.Event, *errors.ErrorTrace) {
	if len(events) == 0 {
		return events, nil
	}

	eventMap := map[types.ID]primitives.Event{}
	for _, event := range events {
		eventMap[event.GetId()] = event
	}

	dbEvents, err := q.getEventEntries(events)
	if err != nil {
		return nil, err.
			Append(errors.LvlWordy, "Could not get cached events").
			Append(errors.LvlPlain, "Database error")
	}

	for _, dbEvent := range dbEvents {
		if event, ok := eventMap[dbEvent.Id]; ok {
			if event.GetColor() == nil {
				event.SetColor(types.ColorFromBytes(dbEvent.Color))
			}
		}
	}

	err = q.insertEvents(events)
	if err != nil {
		return nil, err.
			Append(errors.LvlWordy, "Could not cache events").
			Append(errors.LvlPlain, "Database error")
	}

	return events, nil
}

func (q *Queries) GetEvent(userId types.ID, eventId types.ID) (primitives.Event, *errors.ErrorTrace) {
	decryptionKey, tr := util.GetUserDecryptionKey(q.CommonConfig, userId)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not get event %v", eventId).
			AltStr(errors.LvlBroad, "Could not get event")
	}

	scanner := parsing.NewPgxScanner(q.PrimitivesParser, q)
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

	err := q.Tx.QueryRow(
		q.Context,
		query,
		eventId.UUID(),
		userId.UUID(),
		decryptionKey,
	).Scan(params...)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return nil, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Event %v for user %v not found", eventId, userId).
			AltStr(errors.LvlPlain, "Event not found").
			AltStr(errors.LvlBroad, "Could not get event")
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get event %v for user %v", eventId, userId).
			AltStr(errors.LvlBroad, "Could not get event")
	}

	event, tr := scanner.GetEvent()
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not parse event %v for user %v", eventId, userId).
			AltStr(errors.LvlWordy, "Could not parse event").
			AltStr(errors.LvlBroad, "Could not get event")
	}

	return event, nil
}

func (q *Queries) InsertEvent(event primitives.Event) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
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
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not insert event %v", event.GetName()).
			AltStr(errors.LvlBroad, "Could not add event")
	}

	return nil
}

func (q *Queries) UpdateEvent(event primitives.Event) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		UPDATE events
		SET color = $2, settings = $3
		WHERE id = $1;
		`,
		event.GetId().UUID(),
		event.GetColor().Bytes(),
		event.GetSettings().Bytes(),
	)

	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Event %v not found", event.GetId()).
			AltStr(errors.LvlPlain, "Event not found").
			Append(errors.LvlPlain, "Could not update event %v", event.GetName()).
			AltStr(errors.LvlBroad, "Could not edit event")
	default:
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not update event %v", event.GetId()).
			Append(errors.LvlPlain, "Could not update event %v", event.GetName()).
			AltStr(errors.LvlBroad, "Could not edit event")
	}
}

func (q *Queries) DeleteEvent(userId types.ID, eventId types.ID) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
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

	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Event %v for user %v not found", eventId, userId).
			AltStr(errors.LvlPlain, "Event not found").
			Append(errors.LvlDebug, "Could not delete event %v", eventId).
			AltStr(errors.LvlBroad, "Could not delete event")
	default:
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not delete event %v", eventId).
			AltStr(errors.LvlBroad, "Could not delete event")
	}
}

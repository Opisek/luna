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

func (q *Queries) insertCalendars(cals []primitives.Calendar) *errors.ErrorTrace {
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
			cal.GetSource().GetId(),
			colBytes,
			cal.GetSettings().Bytes(),
		}

		rows = append(rows, row)
	}

	err := util.CopyAndUpdate(
		q.Tx,
		q.Context,
		"calendars",
		[]string{"id", "source", "color", "settings"},
		[]string{"color", "settings"},
		rows,
	)

	if err != nil {
		return err.
			Append(errors.LvlWordy, "Could not insert calendars")
	}

	return nil
}

func (q *Queries) getCalendarEntries(cals []primitives.Calendar) ([]*types.CalendarDatabaseEntry, *errors.ErrorTrace) {
	query := fmt.Sprintf(
		`
		SELECT id, source, color, settings
		FROM calendars
		WHERE id IN (
			%s
		);
		`,
		util.GenerateArgList(1, len(cals)),
	)

	rows, err := q.Tx.Query(
		q.Context,
		query,
		util.JoinIds(cals, func(c primitives.Calendar) types.ID { return c.GetId() })...,
	)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get calendars from the database")
	}
	defer rows.Close()

	entries := []*types.CalendarDatabaseEntry{}
	for rows.Next() {
		entry := &types.CalendarDatabaseEntry{}

		err := rows.Scan(&entry.Id, &entry.Source, &entry.Color, &entry.Settings)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not scan calendar row")
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (q *Queries) ReconcileCalendars(cals []primitives.Calendar) ([]primitives.Calendar, *errors.ErrorTrace) {
	if len(cals) == 0 {
		return cals, nil
	}

	calMap := map[types.ID]primitives.Calendar{}
	for _, cal := range cals {
		calMap[cal.GetId()] = cal
	}

	dbCals, err := q.getCalendarEntries(cals)
	if err != nil {
		return nil, err.
			Append(errors.LvlWordy, "Could not get cached events").
			Append(errors.LvlPlain, "Database error")
	}

	for _, dbCal := range dbCals {
		if cal, ok := calMap[dbCal.Id]; ok {
			if cal.GetColor() == nil {
				cal.SetColor(types.ColorFromBytes(dbCal.Color))
				// TODO: if dbCal.Color == nil, either return some default color, or generate a deterministic random one (e.g. calendar id hash -> hue)
			}
		}
	}

	err = q.insertCalendars(cals)
	if err != nil {
		return nil, err.
			Append(errors.LvlWordy, "Could not cache events").
			Append(errors.LvlPlain, "Database error")
	}

	return cals, nil
}

func (q *Queries) GetCalendar(userId types.ID, calendarId types.ID) (primitives.Calendar, *errors.ErrorTrace) {
	decryptionKey, tr := util.GetUserDecryptionKey(q.CommonConfig, userId)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not get calendar %v", calendarId).
			AltStr(errors.LvlBroad, "Could not get calendar")
	}

	scanner := parsing.NewPgxScanner(q.PrimitivesParser, q)
	scanner.ScheduleCalendar()
	cols, params := scanner.Variables(3)

	query := fmt.Sprintf(
		`
		SELECT %s
		FROM calendars
		JOIN sources ON calendars.source = sources.id
		WHERE calendars.id = $1
		AND sources.userid = $2;
		`,
		cols,
	)

	err := q.Tx.QueryRow(
		q.Context,
		query,
		calendarId.UUID(),
		userId.UUID(),
		decryptionKey,
	).Scan(params...)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return nil, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Calendar %v for user %v not found", calendarId, userId).
			AltStr(errors.LvlPlain, "Calendar not found").
			AltStr(errors.LvlBroad, "Could not get event")
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get calendar %v for user %v", calendarId, userId).
			AltStr(errors.LvlBroad, "Could not get calendar")
	}

	event, tr := scanner.GetCalendar()
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not parse calendar %v for user %v", calendarId, userId).
			AltStr(errors.LvlWordy, "Could not parse calendar").
			AltStr(errors.LvlBroad, "Could not get calendar")
	}

	return event, nil
}

func (q *Queries) InsertCalendar(calendar primitives.Calendar) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		INSERT INTO calendars (id, source, color, settings)
		VALUES ($1, $2, $3, $4);
		`,
		calendar.GetId().UUID(),
		calendar.GetSource().GetId().UUID(),
		calendar.GetColor().Bytes(),
		calendar.GetSettings().Bytes(),
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not insert calendar %v", calendar.GetName()).
			AltStr(errors.LvlBroad, "Could not add calendar")
	}

	return nil
}

func (q *Queries) UpdateCalendar(cal primitives.Calendar) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		UPDATE calendars
		SET color = $1, settings = $2
		WHERE id = $3;`,
		cal.GetColor().Bytes(),
		cal.GetSettings().Bytes(),
		cal.GetId(),
	)

	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Calendar %v not found", cal.GetId()).
			AltStr(errors.LvlPlain, "Calendar not found").
			Append(errors.LvlPlain, "Could not update calendar %v", cal.GetName()).
			AltStr(errors.LvlBroad, "Could not edit calendar")
	default:
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not update calendar %v", cal.GetId()).
			Append(errors.LvlPlain, "Could not update calendar %v", cal.GetName()).
			AltStr(errors.LvlBroad, "Could not edit calendar")
	}
}

func (q *Queries) DeleteCalendar(userId types.ID, calendarId types.ID) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		DELETE FROM calendars
		WHERE id = $1
		AND source IN (
			SELECT id
			FROM sources
			WHERE userid = $2
		);
		`,
		calendarId.UUID(),
		userId.UUID(),
	)

	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Calendar %v for user %v not found", calendarId, userId).
			AltStr(errors.LvlPlain, "Calendar not found").
			Append(errors.LvlDebug, "Could not delete calendar %v", calendarId).
			AltStr(errors.LvlBroad, "Could not delete calendar")
	default:
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not delete calendar %v", calendarId).
			AltStr(errors.LvlBroad, "Could not delete calendar")
	}
}

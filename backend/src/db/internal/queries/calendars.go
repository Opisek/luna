package queries

import (
	"context"
	"fmt"
	"luna-backend/config"

	"luna-backend/db/internal/parsing"
	"luna-backend/db/internal/util"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) insertCalendars(cals []types.Calendar) *errors.ErrorTrace {
	if len(cals) == 0 {
		return nil
	}

	rows := [][]any{}

	for _, cal := range cals {
		row := []any{
			cal.GetId(),
			cal.GetSource().GetId(),
			cal.GetSettings().Bytes(),
		}

		rows = append(rows, row)
	}

	err := util.CopyAndUpdate(
		q.Tx,
		q.Context,
		"calendars",
		"id",
		[]string{"id", "source", "settings"},
		[]string{"settings"},
		rows,
		true,
		"original.source = $1",
		cals[0].GetSource().GetId(),
		true,
		"source",
		"id",
	)

	if err != nil {
		return err.
			Append(errors.LvlWordy, "Could not insert calendars")
	}

	return nil
}

func (q *Queries) getCalendarEntries(cals []types.Calendar) ([]*types.CalendarExtendedDatabaseEntry, *errors.ErrorTrace) {
	query := fmt.Sprintf(
		`
		SELECT id, source, settings, COALESCE(title, '') as title, COALESCE(description, '') as description, color, COALESCE(overridden, false) AS overridden
		FROM calendars
		LEFT OUTER JOIN (
			SELECT calendarid, title, description, color, true AS overridden
			FROM calendar_overrides	
		) AS overrides ON calendars.id = overrides.calendarid
		WHERE id IN (
			%s
		);
		`,
		util.GenerateArgList(1, len(cals)),
	)

	rows, err := q.Tx.Query(
		q.Context,
		query,
		util.JoinIds(cals, func(c types.Calendar) types.ID { return c.GetId() })...,
	)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get calendars from the database")
	}
	defer rows.Close()

	entries := []*types.CalendarExtendedDatabaseEntry{}
	for rows.Next() {
		entry := &types.CalendarExtendedDatabaseEntry{}

		err := rows.Scan(&entry.Id, &entry.Source, &entry.Settings, &entry.Title, &entry.Description, &entry.Color, &entry.Overridden)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not scan calendar row")
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (q *Queries) orderCalendars(cals []types.Calendar) ([]types.Calendar, *errors.ErrorTrace) {
	query := fmt.Sprintf(
		`
		SELECT id
		FROM calendars
		WHERE id IN (
			%s
		)
		ORDER BY source, display_order;
		`,
		util.GenerateArgList(1, len(cals)),
	)

	rows, err := q.Tx.Query(
		q.Context,
		query,
		util.JoinIds(cals, func(c types.Calendar) types.ID { return c.GetId() })...,
	)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get calendars from the database")
	}
	defer rows.Close()

	orderedIds := make([]types.ID, len(cals))
	i := 0
	for rows.Next() {
		err := rows.Scan(&orderedIds[i])
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not scan ordered calendar ID")
		}
		i += 1
	}

	calMap := map[types.ID]types.Calendar{}
	for _, cal := range cals {
		calMap[cal.GetId()] = cal
	}

	orderedCalendars := make([]types.Calendar, len(cals))
	for i, id := range orderedIds {
		cal, exists := calMap[id]
		if !exists {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not retrieve mapped calendar")
		}
		orderedCalendars[i] = cal
	}

	return orderedCalendars, nil
}

func (q *Queries) OverrideCalendars(cals []types.Calendar) ([]types.Calendar, *errors.ErrorTrace) {
	if len(cals) == 0 {
		return cals, nil
	}

	calMap := map[types.ID]types.Calendar{}
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
			if !dbCal.Overridden {
				continue
			}
			cal.SetOverridden(true)
			if dbCal.Title != "" {
				cal.SetName(dbCal.Title)
			}
			if dbCal.Description != "" {
				cal.SetDesc(dbCal.Description)
			}
			if dbCal.Color != nil {
				cal.SetColor(types.ColorFromBytes(dbCal.Color))
			}
		}
	}

	err = q.insertCalendars(cals)
	if err != nil {
		return nil, err.
			Append(errors.LvlWordy, "Could not cache events").
			Append(errors.LvlPlain, "Database error")
	}

	cals, err = q.orderCalendars(cals)
	if err != nil {
		return nil, err.
			Append(errors.LvlWordy, "Could not order calendars").
			Append(errors.LvlPlain, "Database error")
	}

	return cals, nil
}

func (q *Queries) OverrideCalendar(calendar types.Calendar) (types.Calendar, *errors.ErrorTrace) {
	cals, tr := q.OverrideCalendars([]types.Calendar{calendar})
	if tr != nil {
		return nil, tr
	}
	return cals[0], nil
}

func (q *Queries) GetCalendar(userId types.ID, calendarId types.ID, ctx context.Context, config *config.CommonConfig) (types.Calendar, *errors.ErrorTrace) {
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

	event, tr := scanner.GetCalendar(ctx)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not parse calendar %v for user %v", calendarId, userId).
			AltStr(errors.LvlWordy, "Could not parse calendar").
			AltStr(errors.LvlBroad, "Could not get calendar")
	}

	return event, nil
}

func (q *Queries) InsertCalendar(calendar types.Calendar) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		INSERT INTO calendars (id, source, settings, display_order)
		SELECT $1, $2, $3, COALESCE(MAX(display_order) + 1, 0)
		FROM calendars
		WHERE source = $2;
		`,
		calendar.GetId().UUID(),
		calendar.GetSource().GetId().UUID(),
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

func (q *Queries) UpdateCalendar(cal types.Calendar) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		UPDATE calendars
		SET settings = $2
		WHERE id = $1;`,
		cal.GetId(),
		cal.GetSettings().Bytes(),
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

func (q *Queries) UpdateCalendarDisplayOrder(userId types.ID, calendarId types.ID, newIndex uint16) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		WITH moved AS (
			SELECT source AS moved_source, display_order AS old_index, SIGN(display_order - $3) AS direction
			FROM calendars
			WHERE id = $2	
			AND source IN (
				SELECT id
				FROM sources
				WHERE userid = $1
			)
		)
		UPDATE calendars
		SET display_order = CASE
			WHEN id = $2 THEN $3
			ELSE display_order + direction
		END 
		FROM moved
		WHERE source = (SELECT moved_source from moved)
		AND display_order BETWEEN SYMMETRIC $3 AND (SELECT old_index FROM moved);
		`,
		userId.UUID(),
		calendarId.UUID(),
		newIndex,
	)

	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "Calendar %v not found", calendarId).
			AltStr(errors.LvlPlain, "Calendar not found").
			Append(errors.LvlDebug, "Could not reorder calendar %v", calendarId).
			AltStr(errors.LvlBroad, "Could not reorder calendar")
	default:
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not reorder calendar %v", calendarId).
			AltStr(errors.LvlBroad, "Could not reorder calendar")
	}
}

func (q *Queries) DeleteCalendar(userId types.ID, calendarId types.ID) *errors.ErrorTrace {
	var deletedCalendarSourceId types.ID
	var deletedCalendarDisplayOrder int
	err := q.Tx.Conn().QueryRow(
		q.Context,
		`
		DELETE FROM calendars
		WHERE id = $1
		AND source IN (
			SELECT id
			FROM sources
			WHERE userid = $2
		)
		RETURNING source, display_order;
		`,
		calendarId.UUID(),
		userId.UUID(),
	).Scan(&deletedCalendarSourceId, &deletedCalendarDisplayOrder)

	switch err {
	case nil:
		break
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

	// Decrease the display order of the sources's other calendars to fill in the gap
	_, err = q.Tx.Exec(
		q.Context,
		`
		UPDATE calendars
		SET display_order = display_order - 1
		WHERE source = $1 AND display_order > $2;
		`,
		deletedCalendarSourceId,
		deletedCalendarDisplayOrder,
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not update display order of other calendars of source %v", deletedCalendarSourceId).
			Append(errors.LvlDebug, "Could not delete calendar %v", calendarId).
			AltStr(errors.LvlBroad, "Could not delete calendar")
	}

	return nil
}

func (q *Queries) SetCalendarOverrides(calendarId types.ID, name string, desc string, color *types.Color) *errors.ErrorTrace {
	columns := []string{}
	params := []any{calendarId.UUID()}

	if name != "" {
		columns = append(columns, "title")
		params = append(params, name)
	}
	if desc != "" {
		columns = append(columns, "description")
		params = append(params, desc)
	}
	if color != nil {
		columns = append(columns, "color")
		params = append(params, color.Bytes())
	}

	query := fmt.Sprintf(
		`
		INSERT INTO calendar_overrides (calendarid, %s)
		VALUES ($1, %s)
		ON CONFLICT (calendarid) DO UPDATE
		SET %s;
		`,
		strings.Join(columns, ", "),
		util.GenerateArgList(2, len(columns)),
		util.GenerateSetList(2, columns),
	)

	_, err := q.Tx.Exec(
		q.Context,
		query,
		params...,
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not set calendar overrides for %v", calendarId).
			AltStr(errors.LvlWordy, "Could not set calendar overrides for %v", name)
	}

	return nil
}

func (q *Queries) DeleteCalendarOverrides(calendarId types.ID) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		DELETE FROM calendar_overrides
		WHERE calendarid = $1;
		`,
		calendarId.UUID(),
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not delete calendar overrides for %v", calendarId).
			AltStr(errors.LvlWordy, "Could not delete calendar overrides")
	}

	return nil
}

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

func (q *Queries) insertCalendars(cals []primitives.Calendar) error {
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

func (q *Queries) getCalendarEntries(sources []primitives.Source, cals []primitives.Calendar) ([]*tables.CalendarEntry, error) {
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
		context.TODO(),
		query,
		util.JoinIds(cals, func(c primitives.Calendar) types.ID { return c.GetId() })...,
	)

	if err != nil {
		return nil, fmt.Errorf("could not get calendars from database: %v", err)
	}

	defer rows.Close()

	sourceMap := map[types.ID]primitives.Source{}
	for _, source := range sources {
		sourceMap[source.GetId()] = source
	}

	entries := []*tables.CalendarEntry{}
	for rows.Next() {
		entry := &tables.CalendarEntry{}

		err := rows.Scan(&entry.Id, &entry.Source, &entry.Color, &entry.Settings)
		if err != nil {
			return nil, fmt.Errorf("could not scan calendar row: %v", err)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (q *Queries) ReconcileCalendars(sources []primitives.Source, cals []primitives.Calendar) ([]primitives.Calendar, error) {
	if len(cals) == 0 {
		return cals, nil
	}

	calMap := map[types.ID]primitives.Calendar{}
	for _, cal := range cals {
		calMap[cal.GetId()] = cal
	}

	dbCals, err := q.getCalendarEntries(sources, cals)
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

	err = q.insertCalendars(cals)
	if err != nil {
		return nil, fmt.Errorf("could not cache calendars: %v", err)
	}

	return cals, nil
}

func (q *Queries) GetCalendar(userId types.ID, calendarId types.ID) (primitives.Calendar, error) {
	var err error

	decryptionKey, err := util.GetUserDecryptionKey(q.CommonConfig, userId)
	if err != nil {
		return nil, fmt.Errorf("could not get user decryption key: %v", err)
	}

	scanner := parsing.NewPgxScanner()
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

	err = q.Tx.QueryRow(
		context.TODO(),
		query,
		calendarId.UUID(),
		userId.UUID(),
		decryptionKey,
	).Scan(params...)
	if err != nil {
		return nil, fmt.Errorf("could not get calendar: %v", err)
	}

	return scanner.GetCalendar()
}

func (q *Queries) UpdateCalendar(cal primitives.Calendar) error {
	_, err := q.Tx.Exec(
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

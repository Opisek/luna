package parsing

import (
	"fmt"
	"luna-backend/errors"
	"luna-backend/types"
	"reflect"
	"strings"
)

type PgxScanner struct {
	source       *types.SourceDatabaseEntry
	scanSource   bool
	calendar     *types.CalendarDatabaseEntry
	scanCalendar bool
	event        *types.EventDatabaseEntry
	scanEvent    bool

	primitivesParser PrimitivesParser
	queries          types.DatabaseQueries
}

func NewPgxScanner(primitivesParser *PrimitivesParser, queries types.DatabaseQueries) *PgxScanner {
	return &PgxScanner{
		scanSource:       false,
		scanCalendar:     false,
		scanEvent:        false,
		primitivesParser: *primitivesParser,
		queries:          queries,
	}
}

func (s *PgxScanner) ScheduleSource() {
	s.source = &types.SourceDatabaseEntry{}
	s.scanSource = true
}

func (s *PgxScanner) ScheduleCalendar() {
	s.ScheduleSource()
	s.calendar = &types.CalendarDatabaseEntry{}
	s.scanCalendar = true
}

func (s *PgxScanner) ScheduleEvent() {
	s.ScheduleCalendar()
	s.event = &types.EventDatabaseEntry{}
	s.scanEvent = true
}

func (s *PgxScanner) ScanSource() bool {
	return s.scanSource
}

func (s *PgxScanner) Variables(keyPos int) (string, []any) {
	vars := []any{}
	columns := []string{}

	tables := []interface{}{}
	tableNames := []string{}

	if s.scanSource {
		tables = append(tables, s.source)
		tableNames = append(tableNames, "sources")
	}
	if s.scanCalendar {
		tables = append(tables, s.calendar)
		tableNames = append(tableNames, "calendars")
	}
	if s.scanEvent {
		tables = append(tables, s.event)
		tableNames = append(tableNames, "events")
	}

	for i, table := range tables {
		t := reflect.TypeOf(table).Elem()
		v := reflect.ValueOf(table).Elem()

		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)

			ptr := v.Field(j).Addr().Interface()
			column := field.Tag.Get("db")
			encrypt := field.Tag.Get("encrypted")

			vars = append(vars, ptr)

			var colStr string
			if encrypt == "true" {
				colStr = fmt.Sprintf("PGP_SYM_DECRYPT(%s.%s, $%d)", tableNames[i], column, keyPos)
			} else {
				colStr = fmt.Sprintf("%s.%s", tableNames[i], column)
			}
			columns = append(columns, colStr)
		}
	}

	return strings.Join(columns, ", "), vars
}

func (s *PgxScanner) GetSourceEntry() *types.SourceDatabaseEntry {
	return s.source
}

func (s *PgxScanner) GetCalendarEntry() *types.CalendarDatabaseEntry {
	return s.calendar
}

func (s *PgxScanner) GetEventEntry() *types.EventDatabaseEntry {
	return s.event
}

func (s *PgxScanner) GetSource() (types.Source, *errors.ErrorTrace) {
	return s.primitivesParser.ParseSource(s.source)
}

func (s *PgxScanner) GetCalendar() (types.Calendar, *errors.ErrorTrace) {
	source, err := s.GetSource()
	if err != nil {
		return nil, err
	}

	settings, err := s.primitivesParser.ParseCalendarSettings(source.GetType(), s.calendar.Settings)
	if err != nil {
		return nil, err
	}

	calendar, err := source.GetCalendar(settings, s.queries)
	if err != nil {
		return nil, err
	}

	return calendar, nil
}

func (s *PgxScanner) GetEvent() (types.Event, *errors.ErrorTrace) {
	calendar, err := s.GetCalendar()
	if err != nil {
		return nil, err
	}

	settings, err := s.primitivesParser.ParseEventSettings(calendar.GetSource().GetType(), s.event.Settings)
	if err != nil {
		return nil, err
	}

	event, err := calendar.GetEvent(settings, s.queries)
	if err != nil {
		return nil, err
	}

	return event, nil
}

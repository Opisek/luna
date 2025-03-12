package parsing

import (
	"fmt"
	"luna-backend/interface/primitives"
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

func (s *PgxScanner) GetSource() (primitives.Source, error) {
	source, err := s.primitivesParser.ParseSource(s.source)
	if err != nil {
		return nil, fmt.Errorf("could not get source: %v", err)
	}
	return source, nil
}

func (s *PgxScanner) GetCalendar() (primitives.Calendar, error) {
	source, err := s.GetSource()
	if err != nil {
		return nil, fmt.Errorf("could not get calendar: %v", err)
	}

	settings, err := s.primitivesParser.ParseCalendarSettings(source.GetType(), s.calendar.Settings)
	if err != nil {
		return nil, fmt.Errorf("could not get calendar:  %v", err)
	}

	calendar, err := source.GetCalendar(settings, s.queries)
	if err != nil {
		return nil, fmt.Errorf("could not get calendar: %v", err)
	}

	if calendar.GetColor() == nil {
		calendar.SetColor(types.ColorFromBytes(s.calendar.Color))
	}

	return calendar, nil
}

func (s *PgxScanner) GetEvent() (primitives.Event, error) {
	calendar, err := s.GetCalendar()
	if err != nil {
		return nil, fmt.Errorf("could not parse event: %v", err)
	}

	settings, err := s.primitivesParser.ParseEventSettings(calendar.GetSource().GetType(), s.event.Settings)
	if err != nil {
		return nil, fmt.Errorf("could not parse event: %v", err)
	}

	event, err := calendar.GetEvent(settings, s.queries)
	if err != nil {
		return nil, fmt.Errorf("could not parse event: %v", err)
	}

	if event.GetColor() == nil {
		event.SetColor(types.ColorFromBytes(s.event.Color))
	}

	return event, nil
}

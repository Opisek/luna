package ical

import (
	"encoding/json"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"time"

	"github.com/emersion/go-ical"
)

type IcalCalendar struct {
	name         string
	desc         string
	source       *IcalSource
	color        *types.Color
	settings     *IcalCalendarSettings
	icalCalendar *ical.Calendar
}

type IcalCalendarSettings struct {
}

func (settings *IcalCalendarSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genCalId(sourceId types.ID, uid string) types.ID {
	return crypto.DeriveID(sourceId, uid)
}

func (calendar *IcalCalendar) GetId() types.ID {
	// ical files only have a single calendar, so they sometimes don't come with a unique ID
	return genCalId(calendar.source.id, "calendar")
}

func (calendar *IcalCalendar) GetName() string {
	return calendar.name
}

func (calendar *IcalCalendar) GetDesc() string {
	return calendar.desc
}

func (calendar *IcalCalendar) GetSource() primitives.Source {
	return calendar.source
}

func (calendar *IcalCalendar) GetSettings() primitives.CalendarSettings {
	return calendar.settings
}

func (calendar *IcalCalendar) GetColor() *types.Color {
	if calendar.color == nil {
		return types.ColorEmpty
	} else {
		return calendar.color
	}
}

func (calendar *IcalCalendar) SetColor(color *types.Color) {
	calendar.color = color
}

func (calendar *IcalCalendar) GetEvents(start time.Time, end time.Time, q types.DatabaseQueries) ([]primitives.Event, error) {
	return nil, fmt.Errorf("not implemented")
}

func (calendar *IcalCalendar) GetEvent(settings primitives.EventSettings, q types.DatabaseQueries) (primitives.Event, error) {
	return nil, fmt.Errorf("not implemented")
}

/* Ical calendar is read-only */

func (calendar *IcalCalendar) AddEvent(name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (primitives.Event, error) {
	return nil, fmt.Errorf("not supported")
}

func (calendar *IcalCalendar) EditEvent(originalEvent primitives.Event, name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (primitives.Event, error) {
	return nil, fmt.Errorf("not supported")
}

func (calendar *IcalCalendar) DeleteEvent(event primitives.Event, q types.DatabaseQueries) error {
	return fmt.Errorf("not supported")
}

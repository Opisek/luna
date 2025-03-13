package ical

import (
	"encoding/json"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/interface/primitives"
	common "luna-backend/interface/protocols/internal"
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

func (source *IcalSource) calendarFromIcal(rawCalendar *ical.Calendar) (*IcalCalendar, error) {
	name := rawCalendar.Props.Get(ical.PropName)
	if name == nil {
		name = rawCalendar.Props.Get("X-WR-CALNAME")
	}
	if name == nil {
		name = ical.NewProp(ical.PropName)
		name.SetText(source.name)
	}

	desc := rawCalendar.Props.Get(ical.PropDescription)
	if desc == nil {
		desc = rawCalendar.Props.Get("X-WR-CALDESC")
	}
	if desc == nil {
		desc = ical.NewProp(ical.PropDescription)
		desc.SetText("")
	}

	var calColor *types.Color = nil
	colProp := rawCalendar.Props.Get(ical.PropColor)
	if colProp != nil {
		var err error
		calColor, err = types.ParseColor(colProp.Value)
		if err != nil {
			calColor = nil
		}
	}

	settings := &IcalCalendarSettings{}

	calendar := &IcalCalendar{
		name:         common.UnespaceIcalString(name.Value),
		desc:         common.UnespaceIcalString(desc.Value),
		source:       source,
		color:        calColor,
		settings:     settings,
		icalCalendar: rawCalendar,
	}

	return calendar, nil
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
	res := make([]primitives.Event, len(calendar.icalCalendar.Children))

	count := 0
	for _, comp := range calendar.icalCalendar.Children {
		if comp.Name != "VEVENT" {
			continue
		}

		event, err := calendar.eventFromIcal(&comp.Props)
		if err != nil {
			return nil, err
		}

		if !event.GetDate().Recurrence().Repeats() && (event.GetDate().Start().Before(start) || event.GetDate().End().After(end)) {
			continue
		}
		res[count] = event
		count++
	}

	return res[:count], nil
}

func (calendar *IcalCalendar) GetEvent(settings primitives.EventSettings, q types.DatabaseQueries) (primitives.Event, error) {
	icalSettings := settings.(*IcalEventSettings)
	targetUid := icalSettings.Uid

	for _, comp := range calendar.icalCalendar.Children {
		if comp.Name != "VEVENT" {
			continue
		}

		event, err := calendar.eventFromIcal(&comp.Props)
		if err != nil {
			return nil, err
		}

		if event.GetSettings().(*IcalEventSettings).Uid == targetUid {
			return event, nil
		}
	}

	return nil, fmt.Errorf("event not found")
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

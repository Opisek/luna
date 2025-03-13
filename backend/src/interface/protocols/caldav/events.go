package caldav

import (
	"encoding/json"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/interface/primitives"
	common "luna-backend/interface/protocols/internal"
	"luna-backend/types"

	"github.com/emersion/go-webdav/caldav"
)

type CaldavEvent struct {
	name      string
	desc      string
	color     *types.Color
	settings  *CaldavEventSettings
	calendar  *CaldavCalendar
	eventDate *types.EventDate
}

type CaldavEventSettings struct {
	Url      *types.Url             `json:"url"`
	Uid      string                 `json:"uid"`
	rawEvent *caldav.CalendarObject `json:"-"`
}

func (calendar *CaldavCalendar) eventFromCaldav(obj *caldav.CalendarObject, q types.DatabaseQueries) (*CaldavEvent, error) {
	eventIndex := -1
	for i, child := range obj.Data.Children {
		if child.Name == "VEVENT" {
			eventIndex = i
			break
		}
	}
	if eventIndex == -1 {
		return nil, fmt.Errorf("could not find VEVENT in calendar object %v", obj.Path)
	}

	parsedProps, mustUpdate, err := common.ParseIcalEvent(&obj.Data.Children[eventIndex].Props)
	if err != nil {
		return nil, fmt.Errorf("could not parse ical event: %w", err)
	}

	url, err := types.NewUrl(obj.Path)
	if err != nil {
		return nil, fmt.Errorf("could not parse event URL %v: %w", obj.Path, err)
	}

	event := &CaldavEvent{
		name:  parsedProps.Name,
		desc:  parsedProps.Desc,
		color: parsedProps.Color,
		settings: &CaldavEventSettings{
			Url:      url,
			Uid:      parsedProps.Uid,
			rawEvent: obj,
		},
		calendar:  calendar,
		eventDate: parsedProps.EventDate,
	}

	if mustUpdate {
		calendar.EditEvent(event, parsedProps.Name, parsedProps.Desc, parsedProps.Color, parsedProps.EventDate, q)
		// TODO: we might want to catch errors and display them as notifications here
	}

	return event, nil
}

func (settings *CaldavEventSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genEventId(calendarId types.ID, uid string) types.ID {
	return crypto.DeriveID(calendarId, uid)
}

func (event *CaldavEvent) GetId() types.ID {
	return genEventId(event.calendar.GetId(), event.settings.Uid)
}

func (event *CaldavEvent) GetName() string {
	return event.name
}

func (event *CaldavEvent) GetDesc() string {
	return event.desc
}

func (event *CaldavEvent) GetCalendar() primitives.Calendar {
	return event.calendar
}

func (event *CaldavEvent) GetSettings() primitives.EventSettings {
	return event.settings
}

func (event *CaldavEvent) GetColor() *types.Color {
	if event.color == nil {
		return event.calendar.GetColor()
	} else {
		return event.color
	}
}

func (event *CaldavEvent) SetColor(color *types.Color) {
	event.color = color
}

func (event *CaldavEvent) GetDate() *types.EventDate {
	return event.eventDate
}

func (event *CaldavEvent) Clone() primitives.Event {
	return &CaldavEvent{
		name:      event.name,
		desc:      event.desc,
		color:     event.color.Clone(),
		settings:  event.settings,
		calendar:  event.calendar,
		eventDate: event.eventDate.Clone(),
	}
}

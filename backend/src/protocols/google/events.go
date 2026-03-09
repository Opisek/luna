package google

import (
	"encoding/json"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/emersion/go-webdav/caldav"
)

type GoogleEvent struct {
	name       string
	desc       string
	color      *types.Color
	overridden bool
	settings   *GoogleEventSettings
	calendar   *GoogleCalendar
	eventDate  *types.EventDate
}

type GoogleEventSettings struct {
	Url      *types.Url             `json:"url"`
	Uid      string                 `json:"uid"`
	rawEvent *caldav.CalendarObject `json:"-"`
}

func (calendar *GoogleCalendar) eventFromGoogle() (*GoogleEvent, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusNotImplemented)
}

func (settings *GoogleEventSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genEventId(calendarId types.ID, uid string) types.ID {
	return crypto.DeriveID(calendarId, uid)
}

func (event *GoogleEvent) GetId() types.ID {
	return genEventId(event.calendar.GetId(), event.settings.Uid)
}

func (event *GoogleEvent) GetName() string {
	return event.name
}

func (event *GoogleEvent) SetName(name string) {
	event.name = name
}

func (event *GoogleEvent) GetDesc() string {
	return event.desc
}

func (event *GoogleEvent) SetDesc(desc string) {
	event.desc = desc
}

func (event *GoogleEvent) GetCalendar() types.Calendar {
	return event.calendar
}

func (event *GoogleEvent) GetSettings() types.EventSettings {
	return event.settings
}

func (event *GoogleEvent) GetColor() *types.Color {
	if event.color == nil {
		return event.calendar.GetColor()
	} else {
		return event.color
	}
}

func (event *GoogleEvent) SetColor(color *types.Color) {
	event.color = color
}

func (event *GoogleEvent) GetOverridden() bool {
	return event.overridden
}

func (event *GoogleEvent) SetOverridden(overridden bool) {
	event.overridden = overridden
}

func (event *GoogleEvent) GetDate() *types.EventDate {
	return event.eventDate
}

func (event *GoogleEvent) Clone() types.Event {
	return &GoogleEvent{
		name:       event.name,
		desc:       event.desc,
		color:      event.color.Clone(),
		overridden: event.overridden,
		settings:   event.settings,
		calendar:   event.calendar,
		eventDate:  event.eventDate.Clone(),
	}
}

package ical

import (
	"encoding/json"
	"luna-backend/crypto"
	"luna-backend/interface/primitives"
	"luna-backend/types"
)

type IcalEvent struct {
	name      string
	desc      string
	color     *types.Color
	settings  *IcalEventSettings
	calendar  *IcalCalendar
	eventDate *types.EventDate
}

type IcalEventSettings struct {
}

func (settings *IcalEventSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genEventId(calendarId types.ID, uid string) types.ID {
	return crypto.DeriveID(calendarId, uid)
}

func (event *IcalEvent) GetId() types.ID {
	return genEventId(event.calendar.GetId(), "TODO")
}

func (event *IcalEvent) GetName() string {
	return event.name
}

func (event *IcalEvent) GetDesc() string {
	return event.desc
}

func (event *IcalEvent) GetCalendar() primitives.Calendar {
	return event.calendar
}

func (event *IcalEvent) GetSettings() primitives.EventSettings {
	return event.settings
}

func (event *IcalEvent) GetColor() *types.Color {
	if event.color == nil {
		return event.calendar.GetColor()
	} else {
		return event.color
	}
}

func (event *IcalEvent) SetColor(color *types.Color) {
	event.color = color
}

func (event *IcalEvent) GetDate() *types.EventDate {
	return event.eventDate
}

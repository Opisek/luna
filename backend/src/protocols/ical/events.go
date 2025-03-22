package ical

import (
	"encoding/json"
	"luna-backend/crypto"
	"luna-backend/errors"
	common "luna-backend/protocols/internal"
	"luna-backend/types"
	"net/http"

	"github.com/emersion/go-ical"
)

type IcalEvent struct {
	name       string
	desc       string
	color      *types.Color
	overridden bool
	settings   *IcalEventSettings
	calendar   *IcalCalendar
	eventDate  *types.EventDate
}

type IcalEventSettings struct {
	Uid string `json:"uid"`
	//rawEvent *ical.Event `json:"-"`
}

func (calendar *IcalCalendar) eventFromIcal(props *ical.Props) (*IcalEvent, *errors.ErrorTrace) {
	parsedProps, _, err := common.ParseIcalEvent(props)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not parse iCal event")
	}

	event := &IcalEvent{
		name:       parsedProps.Name,
		desc:       parsedProps.Desc,
		color:      parsedProps.Color,
		overridden: false,
		settings: &IcalEventSettings{
			Uid: parsedProps.Uid,
			//rawEvent: icalEvent,
		},
		calendar:  calendar,
		eventDate: parsedProps.EventDate,
	}

	return event, nil
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
	return genEventId(event.calendar.GetId(), event.settings.Uid)
}

func (event *IcalEvent) GetName() string {
	return event.name
}

func (event *IcalEvent) SetName(name string) {
	event.name = name
}

func (event *IcalEvent) GetDesc() string {
	return event.desc
}

func (event *IcalEvent) SetDesc(desc string) {
	event.desc = desc
}

func (event *IcalEvent) GetCalendar() types.Calendar {
	return event.calendar
}

func (event *IcalEvent) GetSettings() types.EventSettings {
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

func (event *IcalEvent) GetOverridden() bool {
	return event.overridden
}

func (event *IcalEvent) SetOverridden(overridden bool) {
	event.overridden = overridden
}

func (event *IcalEvent) GetDate() *types.EventDate {
	return event.eventDate
}

func (event *IcalEvent) Clone() types.Event {
	return &IcalEvent{
		name:       event.name,
		desc:       event.desc,
		color:      event.color.Clone(),
		overridden: event.overridden,
		settings:   event.settings,
		calendar:   event.calendar,
		eventDate:  event.eventDate.Clone(),
	}
}

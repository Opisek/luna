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
	Uid               string `json:"uid"`
	RecurrenceId      string `json:"recurrence_id"`
	IsFirstRecurrence bool   `json:"is_first_recurrence"`
	//rawEvent *ical.Event `json:"-"`
}

func (settings *IcalEventSettings) Clone() *IcalEventSettings {
	return &IcalEventSettings{
		Uid:               settings.Uid,
		RecurrenceId:      settings.RecurrenceId,
		IsFirstRecurrence: settings.IsFirstRecurrence,
	}
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
			Uid:               parsedProps.Uid,
			RecurrenceId:      parsedProps.RecurrenceId,
			IsFirstRecurrence: parsedProps.RecurrenceId == "",
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

func (event *IcalEvent) GetId() types.ID {
	masterEventId := crypto.DeriveID(event.calendar.GetId(), event.settings.Uid)

	if event.settings.RecurrenceId == "" || event.settings.IsFirstRecurrence {
		return masterEventId
	}

	return crypto.DeriveID(masterEventId, event.settings.RecurrenceId)
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
		settings:   event.settings.Clone(),
		calendar:   event.calendar,
		eventDate:  event.eventDate.Clone(),
	}
}

func (event *IcalEvent) SupplyMasterEvent(masterEvent types.Event) {
	event.settings.RecurrenceId = common.CalculateRecurrenceId(event.eventDate.Start(), event.eventDate.AllDay())
	event.settings.IsFirstRecurrence = masterEvent.GetDate().Start().Equal(*event.eventDate.Start())
}

func (event *IcalEvent) GetRecurrenceId() string {
	return event.settings.RecurrenceId
}

func (event *IcalEvent) CanEdit() bool {
	return false
}

func (event *IcalEvent) CanDelete() bool {
	return false
}

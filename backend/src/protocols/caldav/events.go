package caldav

import (
	"encoding/json"
	"luna-backend/crypto"
	"luna-backend/errors"
	common "luna-backend/protocols/internal"
	"luna-backend/types"
	"net/http"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

type CaldavEvent struct {
	name       string
	desc       string
	color      *types.Color
	overridden bool
	settings   *CaldavEventSettings
	calendar   *CaldavCalendar
	eventDate  *types.EventDate
}

type CaldavEventSettings struct {
	Url               *types.Url             `json:"url"`
	Uid               string                 `json:"uid"`
	RecurrenceId      string                 `json:"recurrence_id"`
	IsFirstRecurrence bool                   `json:"is_first_recurrence"`
	rawEvent          *caldav.CalendarObject `json:"-"`
}

func (settings *CaldavEventSettings) Clone() *CaldavEventSettings {
	newUrl := *settings.Url

	return &CaldavEventSettings{
		Url:               &newUrl,
		Uid:               settings.Uid,
		RecurrenceId:      settings.RecurrenceId,
		IsFirstRecurrence: settings.IsFirstRecurrence,
		rawEvent:          settings.rawEvent,
	}
}

func (calendar *CaldavCalendar) eventFromCaldav(obj *caldav.CalendarObject, q types.DatabaseQueries) (*CaldavEvent, *errors.ErrorTrace) {
	eventIndex := -1
	for i, child := range obj.Data.Children {
		if child.Name == "VEVENT" {
			eventIndex = i
			break
		}
	}
	if eventIndex == -1 {
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "could not find VEVENT in calendar object %v", obj.Path)
	}

	parsedProps, mustUpdate, err := common.ParseIcalEvent(&obj.Data.Children[eventIndex].Props)
	if err != nil {
		uid := "unknown"
		if uidProp := obj.Data.Children[eventIndex].Props.Get(ical.PropUID); uidProp != nil {
			uid = uidProp.Value
		}
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not parse iCal event %v", uid)
	}

	url, err := types.NewUrl(obj.Path)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse event URL")
	}

	event := &CaldavEvent{
		name:       parsedProps.Name,
		desc:       parsedProps.Desc,
		color:      parsedProps.Color,
		overridden: false,
		settings: &CaldavEventSettings{
			Url:               url,
			Uid:               parsedProps.Uid,
			RecurrenceId:      parsedProps.RecurrenceId,
			IsFirstRecurrence: parsedProps.RecurrenceId == "",
			rawEvent:          obj,
		},
		calendar:  calendar,
		eventDate: parsedProps.EventDate,
	}

	if mustUpdate {
		calendar.EditEvent(event, parsedProps.Name, parsedProps.Desc, parsedProps.Color, parsedProps.EventDate, false, q)
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

func (event *CaldavEvent) GetId() types.ID {
	masterEventId := crypto.DeriveID(event.calendar.GetId(), event.settings.Uid)

	if event.settings.RecurrenceId == "" || event.settings.IsFirstRecurrence {
		return masterEventId
	}

	return crypto.DeriveID(masterEventId, event.settings.RecurrenceId)
}

func (event *CaldavEvent) GetName() string {
	return event.name
}

func (event *CaldavEvent) SetName(name string) {
	event.name = name
}

func (event *CaldavEvent) GetDesc() string {
	return event.desc
}

func (event *CaldavEvent) SetDesc(desc string) {
	event.desc = desc
}

func (event *CaldavEvent) GetCalendar() types.Calendar {
	return event.calendar
}

func (event *CaldavEvent) GetSettings() types.EventSettings {
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

func (event *CaldavEvent) GetOverridden() bool {
	return event.overridden
}

func (event *CaldavEvent) SetOverridden(overridden bool) {
	event.overridden = overridden
}

func (event *CaldavEvent) GetDate() *types.EventDate {
	return event.eventDate
}

func (event *CaldavEvent) Clone() types.Event {
	return &CaldavEvent{
		name:       event.name,
		desc:       event.desc,
		color:      event.color.Clone(),
		overridden: event.overridden,
		settings:   event.settings.Clone(),
		calendar:   event.calendar,
		eventDate:  event.eventDate.Clone(),
	}
}

func (event *CaldavEvent) SupplyMasterEvent(masterEvent types.Event) {
	event.settings.RecurrenceId = types.SerializeIcalTime(event.eventDate.Start(), event.eventDate.AllDay(), true)
	event.settings.IsFirstRecurrence = masterEvent.GetDate().Start().Equal(*event.eventDate.Start())
}

func (event *CaldavEvent) IsRecurrenceInstance() bool {
	return event.settings.RecurrenceId != ""
}

func (event *CaldavEvent) GetRecurrenceId() string {
	return event.settings.RecurrenceId
}

func (event *CaldavEvent) CanEdit() bool {
	return true
}

func (event *CaldavEvent) CanDelete() bool {
	return true
}

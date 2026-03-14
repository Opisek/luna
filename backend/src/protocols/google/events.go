package google

import (
	"encoding/json"
	"luna-backend/crypto"
	"luna-backend/errors"
	google "luna-backend/protocols/google/internal"
	"luna-backend/types"
	"net/http"
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
	GoogleId     string `json:"google_id"`
	Uid          string `json:"ical_id"`
	RecurrenceId string `json:"recurrence_id"`
}

func (calendar *GoogleCalendar) eventFromGoogle(googleEvent *google.Event, q types.DatabaseQueries) (*GoogleEvent, *errors.ErrorTrace) {
	var col *types.Color
	if googleEvent.ColorId == "" {
		col = calendar.color.Clone()
	} else {
		var tr *errors.ErrorTrace
		col, _, tr = calendar.source.getColorById(googleEvent.ColorId, false, q)
		if tr != nil {
			return nil, tr.
				Append(errors.LvlDebug, "Could not resolve color id %v", googleEvent.ColorId).
				AltStr(errors.LvlWordy, "Could not resolve color id").
				Append(errors.LvlDebug, "Could not parse event %v", googleEvent.Id).
				AltStr(errors.LvlWordy, "Could not parse event")
		}
	}

	settings := &GoogleEventSettings{
		GoogleId:     googleEvent.Id,
		Uid:          googleEvent.IcalUid,
		RecurrenceId: googleEvent.RecurringEventId,
	}

	recurrence, err := types.EventRecurrenceFromLines(googleEvent.Recurrence)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse recurrence %v", googleEvent.Recurrence).
			AltStr(errors.LvlWordy, "Could not recurrence").
			Append(errors.LvlDebug, "Could not parse event %v", googleEvent.Id).
			AltStr(errors.LvlWordy, "Could not parse event")
	}

	startTime, allDay, tr := googleEvent.Start.ParseTimeDefinition()
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse start time").
			Append(errors.LvlDebug, "Could not parse event %v", googleEvent.Id).
			AltStr(errors.LvlWordy, "Could not parse event")
	}
	endTime, _, tr := googleEvent.Start.ParseTimeDefinition()
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse end time").
			Append(errors.LvlDebug, "Could not parse event %v", googleEvent.Id).
			AltStr(errors.LvlWordy, "Could not parse event")
	}

	eventDate := types.NewEventDateFromEndTime(startTime, endTime, allDay, recurrence)

	event := &GoogleEvent{
		name:       googleEvent.Name,
		desc:       googleEvent.Description,
		color:      col,
		overridden: false,
		settings:   settings,
		calendar:   calendar,
		eventDate:  eventDate,
	}

	return event, nil
}

func (settings *GoogleEventSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genEventId(calendarId types.ID, googleId string) types.ID {
	return crypto.DeriveID(calendarId, googleId)
}

func (event *GoogleEvent) GetId() types.ID {
	return genEventId(event.calendar.GetId(), event.settings.GoogleId)
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

func (event *GoogleEvent) CanEdit() bool {
	return !event.eventDate.Recurrence().Repeats()
}

func (event *GoogleEvent) CanDelete() bool {
	return !event.eventDate.Recurrence().Repeats()
}

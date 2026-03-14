package google

import (
	"context"
	"encoding/json"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/net"
	google "luna-backend/protocols/google/internal"
	"luna-backend/types"
	"net/http"
	"time"
)

type GoogleCalendar struct {
	name       string
	desc       string
	color      *types.Color
	overridden bool
	settings   *GoogleCalendarSettings
	source     *GoogleSource
	primary    bool
}

type GoogleCalendarSettings struct {
	GoogleId string `json:"google_id"`
}

func (source *GoogleSource) calendarFromGoogle(calListEntry *google.CalendarListEntry, q types.DatabaseQueries) (*GoogleCalendar, *errors.ErrorTrace) {
	var tr *errors.ErrorTrace

	var col *types.Color
	if calListEntry.BackgroundColor != "" {
		var err error
		col, err = types.ParseColor(calListEntry.BackgroundColor)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not parse color %v", calListEntry.BackgroundColor).
				Append(errors.LvlDebug, "Could not parse calendar %v", calListEntry.Id).
				AltStr(errors.LvlWordy, "Could not parse calendar")
		}
	} else if calListEntry.ColorId != "" {
		col, _, tr = source.getColorById(calListEntry.ColorId, true, q)
		if tr != nil {
			return nil, tr.
				Append(errors.LvlDebug, "Could not parse calendar %v", calListEntry.Id).
				AltStr(errors.LvlWordy, "Could not parse calendar")
		}
	} else {
		col = types.ColorEmpty.Clone()
	}

	settings := &GoogleCalendarSettings{
		GoogleId: calListEntry.Id,
	}

	calendar := &GoogleCalendar{
		name:       calListEntry.Name,
		desc:       calListEntry.Description,
		color:      col,
		overridden: false,
		settings:   settings,
		source:     source,
		primary:    calListEntry.Primary,
	}

	return calendar, nil
}

func (settings *GoogleCalendarSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genCalId(sourceId types.ID, googleId string) types.ID {
	return crypto.DeriveID(sourceId, googleId)
}

func (calendar *GoogleCalendar) GetId() types.ID {
	return genCalId(calendar.source.id, calendar.settings.GoogleId)
}

func (calendar *GoogleCalendar) GetName() string {
	return calendar.name
}

func (calendar *GoogleCalendar) SetName(name string) {
	calendar.name = name
}

func (calendar *GoogleCalendar) GetDesc() string {
	return calendar.desc
}

func (calendar *GoogleCalendar) SetDesc(desc string) {
	calendar.desc = desc
}

func (calendar *GoogleCalendar) GetSource() types.Source {
	return calendar.source
}

func (calendar *GoogleCalendar) GetSettings() types.CalendarSettings {
	return calendar.settings
}

func (calendar *GoogleCalendar) GetColor() *types.Color {
	if calendar.color == nil {
		return types.ColorEmpty
	} else {
		return calendar.color
	}
}

func (calendar *GoogleCalendar) SetColor(color *types.Color) {
	calendar.color = color
}

func (calendar *GoogleCalendar) GetOverridden() bool {
	return calendar.overridden
}

func (calendar *GoogleCalendar) SetOverridden(overridden bool) {
	calendar.overridden = overridden
}

func (calendar *GoogleCalendar) CanEdit() bool {
	return true
}

func (calendar *GoogleCalendar) CanDelete() bool {
	return !calendar.primary
}

func (calendar *GoogleCalendar) CanAddEvents() bool {
	return true
}

func (calendar *GoogleCalendar) GetEvents(start time.Time, end time.Time, q types.DatabaseQueries) ([]types.Event, *errors.ErrorTrace) {
	var res struct {
		Items []*google.Event `json:"items"`
	}

	url := google.ApiUrl().Subpage("calendars", calendar.settings.GoogleId, "events")
	query := url.Query()
	query.Set("timeMin", start.Format(time.RFC3339))
	query.Set("timeMax", end.Format(time.RFC3339))
	url.SetQuery(query)

	tr := net.FetchJson(url, "GET", calendar.source.auth, nil, "", q.GetContext(), &res)
	if tr != nil {
		return nil, tr
	}

	cancelledInstances := make(map[string][]*google.Event)

	result := make([]types.Event, len(res.Items))
	eventCount := 0
	for _, event := range res.Items {
		// Instead of adding an exception to the RRULE,
		// Google calendar returns an additional "copy" of the event with status set to "cancelled".
		// To get actual instances, we would either have to call the instances endpoint for every recurring event,
		// or we resolve the recurrence ourselves (like we already do) but subtract the cancelled instances.
		// We can do this by adding the cancelled instances to the RRULE exceptions.
		if event.Status == "cancelled" {
			cancelled, exists := cancelledInstances[event.RecurringEventId]
			if !exists {
				cancelled = []*google.Event{}
			}
			cancelled = append(cancelled, event)
			cancelledInstances[event.RecurringEventId] = cancelled
			continue
		}

		converted, err := calendar.eventFromGoogle(event, q)
		if err != nil {
			return nil, err.
				Append(errors.LvlBroad, "Could not get events")
		}

		casted := (types.Event)(converted)

		result[eventCount] = casted
		eventCount += 1
	}

	// Add the exception dates
	for _, event := range result[:eventCount] {
		if exceptions, exists := cancelledInstances[event.GetSettings().(*GoogleEventSettings).GoogleId]; exists {
			for _, exception := range exceptions {
				exceptionTime, _, tr := exception.Start.ParseTimeDefinition()
				if tr != nil {
					return nil, tr.
						Append(errors.LvlWordy, "Could not parse exception time").
						Append(errors.LvlDebug, "Could not parse event %v", exception.Id).
						AltStr(errors.LvlWordy, "Could not parse event")
				}
				event.GetDate().Recurrence().AddException(exceptionTime)
			}
		}
	}

	return result[:eventCount], nil
}

func (calendar *GoogleCalendar) GetEvent(settings types.EventSettings, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	googleSettings := settings.(*GoogleEventSettings)

	var res google.Event

	url := google.ApiUrl().Subpage("calendars", calendar.settings.GoogleId, "events", googleSettings.GoogleId)

	tr := net.FetchJson(url, "GET", calendar.source.auth, nil, "", q.GetContext(), &res)
	if tr != nil {
		return nil, tr
	}

	converted, err := calendar.eventFromGoogle(&res, q)
	if err != nil {
		return nil, err.
			Append(errors.LvlBroad, "Could not get event")
	}

	casted := (types.Event)(converted)

	return casted, nil
}

func (calendar *GoogleCalendar) AddEvent(name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	var tr *errors.ErrorTrace

	var colId string

	if color.IsEmpty() {
		colId = ""
	} else {
		colId, tr = calendar.source.getClosestColorId(color, false, q)
		if tr != nil {
			return nil, tr.
				Append(errors.LvlDebug, "Could not add event to calendar %v", calendar.GetId()).
				AltStr(errors.LvlBroad, "Could not add event")
		}
	}

	var start google.TimeDefinition
	var end google.TimeDefinition
	// TODO: timezones
	if date.AllDay() {
		start = google.TimeDefinition{
			Date: date.Start().Format("2006-01-02"),
		}
		end = google.TimeDefinition{
			Date: date.End().Format("2006-01-02"),
		}
	} else {
		start = google.TimeDefinition{
			DateTime: date.Start().Local().Format(time.RFC3339),
		}
		end = google.TimeDefinition{
			DateTime: date.End().Local().Format(time.RFC3339),
		}
	}

	var recurrence []string
	if date.Recurrence().Repeats() {
		recurrence = []string{date.Recurrence().Rule().String()}
	} else {
		recurrence = []string{}
	}

	event := google.Event{
		Name:        name,
		Description: desc,
		ColorId:     colId,
		Start:       start,
		End:         end,
		Recurrence:  recurrence,
	}

	url := google.ApiUrl().Subpage("calendars", calendar.settings.GoogleId, "events")

	var res google.Event

	tr = net.FetchJson(url, "POST", calendar.source.auth, &event, "application/json", q.GetContext(), &res)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not add event to calendar %v", calendar.GetId()).
			AltStr(errors.LvlBroad, "Could not add event")
	}

	converted, err := calendar.eventFromGoogle(&res, q)
	if err != nil {
		return nil, err.
			Append(errors.LvlDebug, "Could not add event to calendar %v", calendar.GetId()).
			AltStr(errors.LvlBroad, "Could not add event")
	}

	casted := (types.Event)(converted)

	return casted, nil
}

func (calendar *GoogleCalendar) EditEvent(originalEvent types.Event, name string, desc string, color *types.Color, date *types.EventDate, _ bool, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	var tr *errors.ErrorTrace

	var colId string

	if color.IsEmpty() {
		colId = ""
	} else {
		colId, tr = calendar.source.getClosestColorId(color, false, q)
		if tr != nil {
			return nil, tr.
				Append(errors.LvlDebug, "Could not edit event %v in calendar %v", originalEvent.GetId(), calendar.GetId()).
				AltStr(errors.LvlBroad, "Could not edit event")
		}
	}

	var start google.TimeDefinition
	var end google.TimeDefinition
	// TODO: timezones
	if date.AllDay() {
		start = google.TimeDefinition{
			Date: date.Start().Format("2006-01-02"),
		}
		end = google.TimeDefinition{
			Date: date.End().Format("2006-01-02"),
		}
	} else {
		start = google.TimeDefinition{
			DateTime: date.Start().Local().Format(time.RFC3339),
		}
		end = google.TimeDefinition{
			DateTime: date.End().Local().Format(time.RFC3339),
		}
	}

	var recurrence []string
	if date.Recurrence().Repeats() {
		recurrence = []string{date.Recurrence().Rule().String()}
	} else {
		recurrence = []string{}
	}

	event := google.Event{
		Name:        name,
		Description: desc,
		ColorId:     colId,
		Start:       start,
		End:         end,
		Recurrence:  recurrence,
	}

	url := google.ApiUrl().Subpage("calendars", calendar.settings.GoogleId, "events", originalEvent.GetSettings().(*GoogleEventSettings).GoogleId)

	var res google.Event

	tr = net.FetchJson(url, "PATCH", calendar.source.auth, &event, "application/json", q.GetContext(), &res)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not edit event %v in calendar %v", originalEvent.GetId(), calendar.GetId()).
			AltStr(errors.LvlBroad, "Could not edit event")
	}

	converted, err := calendar.eventFromGoogle(&res, q)
	if err != nil {
		return nil, err.
			Append(errors.LvlDebug, "Could not edit event %v in calendar %v", originalEvent.GetId(), calendar.GetId()).
			AltStr(errors.LvlBroad, "Could not edit event")
	}

	casted := (types.Event)(converted)

	return casted, nil
}

func (calendar *GoogleCalendar) DeleteEvent(event types.Event, q types.DatabaseQueries) *errors.ErrorTrace {
	googleSettings := event.GetSettings().(*GoogleEventSettings)

	url := google.ApiUrl().Subpage("calendars", calendar.settings.GoogleId, "events", googleSettings.GoogleId)

	_, tr := net.FetchBytes(url, "DELETE", calendar.source.auth, nil, "", "", q.GetContext())
	if tr != nil {
		return tr
	}

	return nil
}

func (calendar *GoogleCalendar) SupplyContext(ctx context.Context) {
	calendar.source.SupplyContext(ctx)
}

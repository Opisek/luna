package caldav

import (
	"context"
	"encoding/json"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/errors"
	supplementary_caldav "luna-backend/protocols/caldav/internal"
	common "luna-backend/protocols/internal"
	"luna-backend/types"
	"net/http"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

type CaldavCalendar struct {
	name       string
	desc       string
	color      *types.Color
	overridden bool
	settings   *CaldavCalendarSettings
	source     *CaldavSource
	client     *caldav.Client
}

type CaldavCalendarSettings struct {
	Url         *types.Url      `json:"url"`
	rawCalendar caldav.Calendar `json:"-"`
}

func (source *CaldavSource) calendarFromCaldav(rawCalendar caldav.Calendar) (*CaldavCalendar, *errors.ErrorTrace) {
	url, err := types.NewUrl(rawCalendar.Path)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse URL %v", rawCalendar.Path).
			Append(errors.LvlWordy, "Could not parse calendar")
	}

	props, tr := supplementary_caldav.PropFind(source.settings.Url, url, []string{"I:calendar-color"}, source.auth, source.ctx)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse calendar")
	}

	var color *types.Color = nil
	colProp, exists := props["I:calendar-color"]
	if exists && colProp.Found {
		color, err = types.ParseColor(colProp.Value)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not parse calendar color %v", colProp.Value).
				Append(errors.LvlWordy, "Could not parse calendar")
		}
	}

	settings := &CaldavCalendarSettings{
		Url:         url,
		rawCalendar: rawCalendar,
	}

	calendar := &CaldavCalendar{
		name:       rawCalendar.Name,
		desc:       rawCalendar.Description,
		color:      color,
		overridden: false,
		settings:   settings,
		source:     source,
		client:     source.client,
	}

	return calendar, nil
}

func (settings *CaldavCalendarSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genCalId(sourceId types.ID, path string) types.ID {
	return crypto.DeriveID(sourceId, path)
}

func (calendar *CaldavCalendar) GetId() types.ID {
	return genCalId(calendar.source.id, calendar.settings.Url.Path)
}

func (calendar *CaldavCalendar) GetName() string {
	return calendar.name
}

func (calendar *CaldavCalendar) SetName(name string) {
	calendar.name = name
}

func (calendar *CaldavCalendar) GetDesc() string {
	return calendar.desc
}

func (calendar *CaldavCalendar) SetDesc(desc string) {
	calendar.desc = desc
}

func (calendar *CaldavCalendar) GetSource() types.Source {
	return calendar.source
}

func (calendar *CaldavCalendar) GetSettings() types.CalendarSettings {
	return calendar.settings
}

func (calendar *CaldavCalendar) GetColor() *types.Color {
	if calendar.color == nil {
		return types.ColorEmpty
	} else {
		return calendar.color
	}
}

func (calendar *CaldavCalendar) SetColor(color *types.Color) {
	calendar.color = color
}

func (calendar *CaldavCalendar) GetOverridden() bool {
	return calendar.overridden
}

func (calendar *CaldavCalendar) SetOverridden(overridden bool) {
	calendar.overridden = overridden
}

func (calendar *CaldavCalendar) CanEdit() bool {
	return true
}

func (calendar *CaldavCalendar) CanDelete() bool {
	return true
}

func (calendar *CaldavCalendar) CanAddEvents() bool {
	return true
}

func (calendar *CaldavCalendar) convertEvent(event *caldav.CalendarObject, q types.DatabaseQueries) ([]types.Event, *errors.ErrorTrace) {
	convertedEvents, err := calendar.eventsFromCaldav(event, q)
	if err != nil {
		return nil, err.
			Append(errors.LvlDebug, "Could not convert calendar %v", event.Path).
			AltStr(errors.LvlWordy, "Could not convert calendar")
	}

	castedEvents := make([]types.Event, len(convertedEvents))
	for i, event := range convertedEvents {
		castedEvents[i] = event
	}

	return castedEvents, nil
}

func (calendar *CaldavCalendar) getEvents(query *caldav.CalendarQuery, q types.DatabaseQueries) ([]types.Event, *errors.ErrorTrace) {
	client, tr := calendar.source.getClient()
	if tr != nil {
		return nil, tr.
			Append(errors.LvlBroad, "Could not get events")
	}

	events, err := client.QueryCalendar(q.GetContext(), calendar.settings.Url.String(), query)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not get events")
	}

	masterEvents := make(map[string]int)
	masterEventIndices := make(map[int]bool)

	convertedEvents := make([]types.Event, 0, len(events))
	for _, event := range events {
		events, tr := calendar.convertEvent(&event, q)
		if tr != nil {
			return nil, tr.
				Append(errors.LvlBroad, "Could not get events")
		}

		for _, event := range events {
			eventSettings := event.GetSettings().(*CaldavEventSettings)

			if len(eventSettings.RecurrenceId) == 0 {
				i := len(convertedEvents)
				masterEvents[eventSettings.Uid] = i
				masterEventIndices[i] = true
			}

			convertedEvents = append(convertedEvents, event)
		}
	}

	// Internally note all the modified recurrence instances for each master event so that we don't expand these later
	for i, event := range convertedEvents {
		if masterEventIndices[i] {
			continue
		}
		eventSettings := event.GetSettings().(*CaldavEventSettings)
		if masterEvent, exists := masterEvents[eventSettings.Uid]; exists {
			convertedEvents[masterEvent].GetDate().Recurrence().MarkModification(common.ExtractDateFromRecurrenceId(event))
			event.SetParent(convertedEvents[masterEvent])
		}
	}

	return convertedEvents, nil
}

func (calendar *CaldavCalendar) getEventsByFilters(filters []caldav.CompFilter, q types.DatabaseQueries) ([]types.Event, *errors.ErrorTrace) {
	return calendar.getEvents(&caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name: "VCALENDAR",
			Comps: []caldav.CalendarCompRequest{{
				Name: "VEVENT",
				Props: []string{
					ical.PropSummary,
					ical.PropUID,
					ical.PropDateTimeStart,
					ical.PropDateTimeEnd,
					ical.PropDateTimeEnd,
					ical.PropRecurrenceID,
				},
			}},
		},
		CompFilter: caldav.CompFilter{
			Name:  "VCALENDAR",
			Comps: filters,
		},
	}, q)
}

func (calendar *CaldavCalendar) GetEvents(start time.Time, end time.Time, q types.DatabaseQueries) ([]types.Event, *errors.ErrorTrace) {
	return calendar.getEventsByFilters([]caldav.CompFilter{
		{
			Name:  "VEVENT",
			Start: start,
			End:   end,
		},
	}, q)
}

func (calendar *CaldavCalendar) GetEvent(settings types.EventSettings, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	caldavSettings := settings.(*CaldavEventSettings)

	// Fetch event directly by path
	obj, err := calendar.client.GetCalendarObject(q.GetContext(), caldavSettings.Url.Path)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not get event")
	}

	events, tr := calendar.convertEvent(obj, q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlBroad, "Could not get event")
	}

	// Find the referred recurrence instance
	var master types.Event
	for _, event := range events {
		currentRecurrenceId := event.GetSettings().(*CaldavEventSettings).RecurrenceId
		if currentRecurrenceId == caldavSettings.RecurrenceId {
			return event, nil
		}
		if len(currentRecurrenceId) == 0 {
			master = event
		}
	}

	// If the specific recurrence ID was not found, but the master event is present, expand the recurrence to a matching id
	if master != nil {
		// Since recurrence IDs refer to the start date of the event instance, we can easily compute bounds for the expansion
		parsedTime, err := types.ParseIcalTimestampAtLocation(caldavSettings.RecurrenceId, master.GetDate().Timezone())
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Invalid recurrence ID %s for event %s", caldavSettings.RecurrenceId, master.GetId()).
				AltStr(errors.LvlWordy, "Invalid recurrence ID").
				Append(errors.LvlBroad, "Could not get event")
		}

		// Search for the instance using its recurrence ID
		end := parsedTime.Add(24 * time.Hour)
		expanded, tr := types.ExpandRecurrence(master, parsedTime, &end)
		if tr != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				AltStr(errors.LvlWordy, "Could not expand recurrence to derive event instance").
				Append(errors.LvlBroad, "Could not get event")
		}

		// Check if the expansion includes the instance
		for _, event := range expanded {
			if event.GetSettings().(*CaldavEventSettings).RecurrenceId == caldavSettings.RecurrenceId {
				return event, nil
			}
		}
	}

	// Error if not found
	return nil, errors.New().Status(http.StatusInternalServerError).
		Append(errors.LvlDebug, "The returned collection from %s did not contain a matching recurrence ID %s", caldavSettings.Url.Path, caldavSettings.RecurrenceId).
		AltStr(errors.LvlWordy, "The returned collection did not contain a matching recurrence ID").
		Append(errors.LvlBroad, "Could not get event")
}

func setEventProps(cal *ical.Calendar, id string, name *string, desc *string, color *types.Color, date *types.EventDate, recurrenceId string) *errors.ErrorTrace {
	var event *ical.Event = nil
	for _, child := range cal.Children {
		if child.Name != "VEVENT" {
			continue
		}

		currentRecurrenceId := child.Props.Get(ical.PropRecurrenceID)

		if (currentRecurrenceId == nil && len(recurrenceId) != 0) || (currentRecurrenceId != nil && currentRecurrenceId.Value != recurrenceId) {
			continue
		}

		event = ical.NewEvent()
		event.Component = child
		break
	}
	if event == nil {
		event = ical.NewEvent()
		cal.Children = append(cal.Children, event.Component)
	}

	event.Props.SetText(ical.PropUID, id)

	if name != nil {
		event.Props.SetText(ical.PropSummary, common.EscapeIcalString(*name))
	}

	if desc != nil {
		if len(*desc) != 0 {
			event.Props.SetText(ical.PropDescription, common.EscapeIcalString(*desc))
		} else {
			event.Props.Del(ical.PropDescription)
		}
	}

	if color.IsEmpty() {
		event.Props.Del(ical.PropColor)
		event.Props.Del(common.PropColor)
		event.Props.Del(common.PropLastColorName)
	} else {
		colorName, exact := types.ColorToName(color)

		// According to the specification, the "COLOR" property must be a named CSS color.
		// To ensure compatibility, we map colors to the closest named CSS color for other clients,
		// and use a custom property for the exac color displayed in Luna.

		event.Props.SetText(ical.PropColor, colorName)
		if exact {
			event.Props.Del(common.PropColor)
			event.Props.Del(common.PropLastColorName)
		} else {
			event.Props.SetText(common.PropColor, color.String())
			// To detect when the color is changed by another client, we store the last color name in a custom property.
			event.Props.SetText(common.PropLastColorName, colorName)
		}
	}

	if date.AllDay() {
		event.Props.SetDate(ical.PropDateTimeStart, *date.Start())
	} else {
		event.Props.SetDateTime(ical.PropDateTimeStart, *date.Start())
	}
	if date.SpecifyDuration() {
		// TODO: figure this out

		return errors.New().Status(http.StatusNotImplemented)
		//event.Props.SetText(ical.PropDuration, *date.Duration())
	} else {
		if date.AllDay() {
			event.Props.SetDate(ical.PropDateTimeEnd, *date.End())
		} else {
			event.Props.SetDateTime(ical.PropDateTimeEnd, *date.End())
		}
		event.Props.Del(ical.PropDuration)
	}

	if recurrenceProps := types.EventRecurrenceToIcal(date.Recurrence()); recurrenceProps != nil {
		if rruleProp := recurrenceProps.Get(ical.PropRecurrenceRule); rruleProp != nil {
			event.Props.Set(rruleProp)
		} else {
			event.Props.Del(ical.PropRecurrenceRule)
		}

		event.Props.Del(ical.PropRecurrenceDates)
		if rdateProps := recurrenceProps.Values(ical.PropRecurrenceDates); len(rdateProps) != 0 {
			for _, rdateProp := range rdateProps {
				event.Props.Add(&rdateProp)
			}
		}

		event.Props.Del(ical.PropExceptionDates)
		if exdateProps := recurrenceProps.Values(ical.PropExceptionDates); len(exdateProps) != 0 {
			for _, exdateProp := range exdateProps {
				event.Props.Add(&exdateProp)
			}
		}
	}

	if len(recurrenceId) != 0 {
		parsedRecurrenceId, err := types.ParseIcalTimestampAtLocation(recurrenceId, date.Timezone())
		if err != nil {
			return errors.New().Status(http.StatusInternalServerError).
				Append(errors.LvlDebug, "Could not parse recurrence ID %s", recurrenceId).
				AltStr(errors.LvlWordy, "Could not parse recurrence ID").
				Append(errors.LvlBroad, "Could not edit event")
		}
		event.Props.SetDateTime(ical.PropRecurrenceID, *parsedRecurrenceId)
	}

	timestamp := time.Now()
	event.Props.SetDateTime(ical.PropDateTimeStamp, timestamp)
	//event.Props.SetDateTime(util.PropTimestamp, timestamp)

	cal.Props.SetText(ical.PropProductID, "Luna")
	cal.Props.SetText(ical.PropVersion, "0.1.0") // TODO: access version from CommonConfig

	return nil
}

func (calendar *CaldavCalendar) AddEvent(name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	id := types.RandomId()
	cal := ical.NewCalendar()

	tr := setEventProps(cal, id.String(), &name, &desc, color, date, "")
	if tr != nil {
		return nil, tr.Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Could not set iCal properties").
			AltStr(errors.LvlPlain, "Malformed settings").
			Append(errors.LvlBroad, "Could not add event")
	}

	path := fmt.Sprintf("%v%v.ics", calendar.settings.Url.Path, id.String())

	_, err := calendar.client.PutCalendarObject(q.GetContext(), path, cal)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not add event")
	}

	obj, err := calendar.client.GetCalendarObject(q.GetContext(), path)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlWordy, "Could not get finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	finishedEvent, tr := calendar.eventsFromCaldav(obj, q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	if len(finishedEvent) == 0 {
		return nil, tr.Status(http.StatusInternalServerError).
			AltStr(errors.LvlWordy, "Could not find finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	return finishedEvent[0], nil
}

func (calendar *CaldavCalendar) EditEvent(originalEvent types.Event, name *string, desc *string, color *types.Color, date *types.EventDate, _ bool, _ string, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	originalCaldavEvent := originalEvent.(*CaldavEvent)
	originalCaldavSettings := originalCaldavEvent.GetSettings().(*CaldavEventSettings)
	uid := originalCaldavSettings.Uid
	recurrenceId := originalCaldavSettings.RecurrenceId
	originalRawEvent := originalCaldavEvent.settings.rawEvent
	cal := originalRawEvent.Data

	tr := setEventProps(cal, uid, name, desc, color, date, recurrenceId)
	if tr != nil {
		return nil, tr.Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Could not set iCal properties").
			AltStr(errors.LvlPlain, "Malformed settings").
			Append(errors.LvlBroad, "Could not add event")
	}

	_, err := calendar.client.PutCalendarObject(q.GetContext(), originalRawEvent.Path, cal)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlWordy, "Could not edit event").
			AltStr(errors.LvlBroad, "Could not edit event")
	}

	obj, err := calendar.client.GetCalendarObject(q.GetContext(), originalRawEvent.Path)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlWordy, "Could not get finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	finishedEvent, tr := calendar.eventsFromCaldav(obj, q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	if len(finishedEvent) == 0 {
		return nil, tr.Status(http.StatusInternalServerError).
			AltStr(errors.LvlWordy, "Could not find finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	return finishedEvent[0], nil
}

func (calendar *CaldavCalendar) DeleteEvent(event types.Event, q types.DatabaseQueries) *errors.ErrorTrace {
	settings := event.GetSettings().(*CaldavEventSettings)

	err := calendar.client.RemoveAll(q.GetContext(), settings.Url.Path)
	if err != nil {
		return errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "event", "CalDAV event").
			Append(errors.LvlBroad, "Could not delete event")
	}

	return nil
}

func (calendar *CaldavCalendar) SupplyContext(ctx context.Context) {
	calendar.source.SupplyContext(ctx)
}

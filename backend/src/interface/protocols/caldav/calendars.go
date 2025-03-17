package caldav

import (
	"encoding/json"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/interface/primitives"
	common "luna-backend/interface/protocols/internal"
	"luna-backend/types"
	"net/http"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

type CaldavCalendar struct {
	name     string
	desc     string
	source   *CaldavSource
	color    *types.Color
	settings *CaldavCalendarSettings
	client   *caldav.Client
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

	settings := &CaldavCalendarSettings{
		Url:         url,
		rawCalendar: rawCalendar,
	}

	calendar := &CaldavCalendar{
		name:     rawCalendar.Name,
		desc:     rawCalendar.Description,
		source:   source,
		color:    nil,
		settings: settings,
		client:   source.settings.client,
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

func (calendar *CaldavCalendar) GetDesc() string {
	return calendar.desc
}

func (calendar *CaldavCalendar) GetSource() primitives.Source {
	return calendar.source
}

func (calendar *CaldavCalendar) GetSettings() primitives.CalendarSettings {
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

func (calendar *CaldavCalendar) convertEvent(event *caldav.CalendarObject, q types.DatabaseQueries) (primitives.Event, *errors.ErrorTrace) {
	convertedEvent, err := calendar.eventFromCaldav(event, q)
	if err != nil {
		return nil, err.
			Append(errors.LvlDebug, "Could not convert calendar %v", event.Path).
			AltStr(errors.LvlWordy, "Could not convert calendar")
	}

	castedEvent := (primitives.Event)(convertedEvent)

	return castedEvent, nil
}

func (calendar *CaldavCalendar) getEvents(query *caldav.CalendarQuery, q types.DatabaseQueries) ([]primitives.Event, *errors.ErrorTrace) {
	client, tr := calendar.source.getClient()
	if tr != nil {
		return nil, tr.
			Append(errors.LvlBroad, "Could not get events")
	}

	events, err := client.QueryCalendar(q.GetContext(), calendar.settings.Url.String(), query)
	if err != nil {
		return nil, errors.InterpretRemoteError(err, "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not get events")
	}

	convertedEvents := make([]primitives.Event, len(events))
	for i, event := range events {
		convertedEvents[i], tr = calendar.convertEvent(&event, q)
		if tr != nil {
			return nil, tr.
				Append(errors.LvlBroad, "Could not get events")
		}
	}

	return convertedEvents, nil
}

func (calendar *CaldavCalendar) GetEvents(start time.Time, end time.Time, q types.DatabaseQueries) ([]primitives.Event, *errors.ErrorTrace) {
	return calendar.getEvents(&caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name: "VCALENDAR",
			Comps: []caldav.CalendarCompRequest{{
				Name: "VEVENT",
				Props: []string{
					"SUMMARY",
					"UID",
					"DTSTART",
					"DTEND",
					"DURATION",
				},
			}},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{{
				Name:  "VEVENT",
				Start: start,
				End:   end,
			}},
		},
	}, q)
}

func (calendar *CaldavCalendar) GetEvent(settings primitives.EventSettings, q types.DatabaseQueries) (primitives.Event, *errors.ErrorTrace) {
	caldavSettings := settings.(*CaldavEventSettings)

	obj, err := calendar.client.GetCalendarObject(q.GetContext(), caldavSettings.Url.Path)
	if err != nil {
		return nil, errors.InterpretRemoteError(err, "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not get event")
	}

	cal, tr := calendar.convertEvent(obj, q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlBroad, "Could not get event")
	}

	return cal, nil
}

func setEventProps(cal *ical.Calendar, id string, name string, desc string, color *types.Color, date *types.EventDate) *errors.ErrorTrace {
	var event *ical.Event = nil
	for _, child := range cal.Children {
		if child.Name == "VEVENT" {
			event = ical.NewEvent()
			event.Component = child
			break
		}
	}
	if event == nil {
		event = ical.NewEvent()
		cal.Children = append(cal.Children, event.Component)
	}

	event.Props.SetText(ical.PropUID, id)

	event.Props.SetText(ical.PropSummary, common.EscapeIcalString(name))

	if desc != "" {
		event.Props.SetText(ical.PropDescription, common.EscapeIcalString(desc))
	} else {
		event.Props.Del(ical.PropDescription)
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

	timestamp := time.Now()
	event.Props.SetDateTime(ical.PropDateTimeStamp, timestamp)
	//event.Props.SetDateTime(util.PropTimestamp, timestamp)

	cal.Props.SetText(ical.PropProductID, "Luna")
	cal.Props.SetText(ical.PropVersion, "0.1.0") // TODO: access version from CommonConfig

	return nil
}

func (calendar *CaldavCalendar) AddEvent(name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (primitives.Event, *errors.ErrorTrace) {
	id := types.RandomId()
	cal := ical.NewCalendar()

	tr := setEventProps(cal, id.String(), name, desc, color, date)
	if tr != nil {
		return nil, tr.Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Could not set iCal properties").
			AltStr(errors.LvlPlain, "Malformed settings").
			Append(errors.LvlBroad, "Could not add event")
	}

	path := fmt.Sprintf("%v%v.ics", calendar.settings.Url.Path, id.String())

	_, err := calendar.client.PutCalendarObject(q.GetContext(), path, cal)
	if err != nil {
		return nil, errors.InterpretRemoteError(err, "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not add event")
	}

	obj, err := calendar.client.GetCalendarObject(q.GetContext(), path)
	if err != nil {
		return nil, errors.InterpretRemoteError(err, "calendar", "CalDAV calendar").
			Append(errors.LvlWordy, "Could not get finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	finishedEvent, tr := calendar.eventFromCaldav(obj, q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	return finishedEvent, nil
}

func (calendar *CaldavCalendar) EditEvent(originalEvent primitives.Event, name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (primitives.Event, *errors.ErrorTrace) {
	originalCaldavEvent := originalEvent.(*CaldavEvent)
	uid := originalCaldavEvent.GetSettings().(*CaldavEventSettings).Uid
	originalRawEvent := originalCaldavEvent.settings.rawEvent
	cal := originalRawEvent.Data

	tr := setEventProps(cal, uid, name, desc, color, date)
	if tr != nil {
		return nil, tr.Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Could not set iCal properties").
			AltStr(errors.LvlPlain, "Malformed settings").
			Append(errors.LvlBroad, "Could not add event")
	}

	_, err := calendar.client.PutCalendarObject(q.GetContext(), originalRawEvent.Path, cal)
	if err != nil {
		return nil, errors.InterpretRemoteError(err, "calendar", "CalDAV calendar").
			Append(errors.LvlWordy, "Could not edit event").
			AltStr(errors.LvlBroad, "Could not edit event")
	}

	obj, err := calendar.client.GetCalendarObject(q.GetContext(), originalRawEvent.Path)
	if err != nil {
		return nil, errors.InterpretRemoteError(err, "calendar", "CalDAV calendar").
			Append(errors.LvlWordy, "Could not get finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	finishedEvent, tr := calendar.eventFromCaldav(obj, q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	return finishedEvent, nil
}

func (calendar *CaldavCalendar) DeleteEvent(event primitives.Event, q types.DatabaseQueries) *errors.ErrorTrace {
	settings := event.GetSettings().(*CaldavEventSettings)

	err := calendar.client.RemoveAll(q.GetContext(), settings.Url.Path)
	if err != nil {
		return errors.InterpretRemoteError(err, "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not delete event")
	}

	return nil
}

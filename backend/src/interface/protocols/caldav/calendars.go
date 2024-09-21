package caldav

import (
	"context"
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/crypto"
	"luna-backend/interface/primitives"
	"luna-backend/types"
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
	auth     auth.AuthMethod
	client   *caldav.Client
}

type CaldavCalendarSettings struct {
	Url         *types.Url      `json:"url"`
	rawCalendar caldav.Calendar `json:"-"`
}

func (source *CaldavSource) calendarFromCaldav(rawCalendar caldav.Calendar) (*CaldavCalendar, error) {
	url, err := types.NewUrl(rawCalendar.Path)
	if err != nil {
		return nil, fmt.Errorf("could not parse calendar URL %v: %w", rawCalendar.Path, err)
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
		client:   source.client,
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

func (calendar *CaldavCalendar) GetAuth() auth.AuthMethod {
	return calendar.auth
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

func (calendar *CaldavCalendar) convertEvent(event *caldav.CalendarObject) (primitives.Event, error) {
	convertedEvent, err := eventFromCaldav(calendar, event)
	if err != nil {
		return nil, fmt.Errorf("could not convert event %v: %w", event.Path, err)
	}

	castedEvent := (primitives.Event)(convertedEvent)

	return castedEvent, nil
}

func (calendar *CaldavCalendar) getEvents(query *caldav.CalendarQuery) ([]primitives.Event, error) {
	client, err := calendar.source.getClient()
	if err != nil {
		return nil, fmt.Errorf("could not get caldav client: %w", err)
	}

	events, err := client.QueryCalendar(context.TODO(), calendar.settings.Url.String(), query)
	if err != nil {
		return nil, fmt.Errorf("could not query calendar: %w", err)
	}

	convertedEvents := make([]primitives.Event, len(events))
	for i, event := range events {
		convertedEvents[i], err = calendar.convertEvent(&event)
		if err != nil {
			return nil, err
		}
	}

	return convertedEvents, nil
}

func (calendar *CaldavCalendar) GetEvents(start time.Time, end time.Time) ([]primitives.Event, error) {
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
	})
}

func (calendar *CaldavCalendar) GetEvent(settings primitives.EventSettings) (primitives.Event, error) {
	caldavSettings := settings.(*CaldavEventSettings)

	obj, err := calendar.client.GetCalendarObject(context.TODO(), caldavSettings.Url.Path)
	if err != nil {
		return nil, fmt.Errorf("could not get event: %w", err)
	}

	cal, err := calendar.convertEvent(obj)
	if err != nil {
		return nil, fmt.Errorf("could not get event: %w", err)
	}

	return cal, nil
}

func (calendar *CaldavCalendar) AddEvent(name string, desc string, color *types.Color, date *types.EventDate) (primitives.Event, error) {
	id := types.RandomId()

	event := ical.NewEvent()

	event.Props.SetText(ical.PropUID, id.String())

	event.Props.SetText(ical.PropSummary, name)

	if desc != "" {
		event.Props.SetText(ical.PropDescription, desc)
	} else {
		event.Props.Del(ical.PropDescription)
	}

	event.Props.SetDateTime(ical.PropDateTimeStart, *date.Start())
	if date.SpecifyDuration() {
		// TODO: figure this out
		return nil, fmt.Errorf("not implemented")
		//event.Props.SetText(ical.PropDuration, *date.Duration())
	} else {
		event.Props.SetDateTime(ical.PropDateTimeEnd, *date.End())
		event.Props.Del(ical.PropDuration)
	}
	event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())

	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropProductID, "Luna")
	cal.Props.SetText(ical.PropVersion, "0.1.0") // TODO: access version from CommonConfig
	cal.Children = append(cal.Children, event.Component)

	path := fmt.Sprintf("%v%v.ics", calendar.settings.Url.Path, id.String())

	_, err := calendar.client.PutCalendarObject(context.TODO(), path, cal)
	if err != nil {
		return nil, fmt.Errorf("could not add event: %w", err)
	}

	obj, err := calendar.client.GetCalendarObject(context.TODO(), path)
	if err != nil {
		return nil, fmt.Errorf("could not get finished event: %w", err)
	}

	finishedEvent, err := eventFromCaldav(calendar, obj)
	if err != nil {
		return nil, fmt.Errorf("could not parse finished event: %w", err)
	}

	return finishedEvent, nil
}

func (calendar *CaldavCalendar) DeleteEvent(settings primitives.EventSettings) error {
	caldavSettings := settings.(*CaldavEventSettings)

	err := calendar.client.RemoveAll(context.TODO(), caldavSettings.Url.Path)
	if err != nil {
		return fmt.Errorf("could not delete event: %w", err)
	}

	return nil
}

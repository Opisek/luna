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
	Url *types.Url `json:"url"`
}

func (source *CaldavSource) calendarFromCaldav(rawCalendar caldav.Calendar) (*CaldavCalendar, error) {
	url, err := types.NewUrl(rawCalendar.Path)
	if err != nil {
		return nil, fmt.Errorf("could not parse calendar URL %v: %w", rawCalendar.Path, err)
	}

	settings := &CaldavCalendarSettings{
		Url: url,
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

func (calendar *CaldavCalendar) GetId() types.ID {
	return crypto.DeriveID(calendar.source.id, calendar.settings.Url.Path)
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

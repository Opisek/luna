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
		return nil, fmt.Errorf("could not parse calendar URL: %w", err)
	}

	settings := &CaldavCalendarSettings{
		Url: url,
	}

	calendar := &CaldavCalendar{
		name:     rawCalendar.Name,
		desc:     rawCalendar.Description,
		source:   source,
		settings: settings,
		client:   source.client,
	}

	return calendar, nil
}

func (settings *CaldavCalendarSettings) GetBytes() []byte {
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

func (calendar *CaldavCalendar) GetSource() types.ID {
	return calendar.source.id
}

func (calendar *CaldavCalendar) GetAuth() auth.AuthMethod {
	return calendar.auth
}

func (calendar *CaldavCalendar) GetSettings() primitives.CalendarSettings {
	return calendar.settings
}

func (calendar *CaldavCalendar) GetColor() *types.Color {
	return types.ColorFromVals(50, 50, 50)
}

func (calendar *CaldavCalendar) GetEvents(start time.Time, end time.Time) ([]primitives.Event, error) {
	client, err := calendar.source.getClient()
	if err != nil {
		return nil, fmt.Errorf("could not get caldav client: %w", err)
	}

	query := caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name:     "VEVENT",
			AllProps: true,
			AllComps: true,
		},
		CompFilter: caldav.CompFilter{},
	}

	events, err := client.QueryCalendar(context.TODO(), calendar.settings.Url.String(), &query)
	if err != nil {
		return nil, fmt.Errorf("could not query calendar: %w", err)
	}

	fmt.Printf("events: %v\n", events)

	return nil, nil
}

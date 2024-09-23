package caldav

import (
	"context"
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/interface/primitives"
	"luna-backend/types"

	"github.com/emersion/go-webdav/caldav"
)

type CaldavSource struct {
	id       types.ID
	name     string
	settings *CaldavSourceSettings
	auth     auth.AuthMethod
	client   *caldav.Client
}

type CaldavSourceSettings struct {
	Url *types.Url `json:"url"`
}

func (settings *CaldavSourceSettings) GetBytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (source *CaldavSource) GetType() string {
	return types.SourceCaldav
}

func (source *CaldavSource) GetId() types.ID {
	return source.id
}

func (source *CaldavSource) GetName() string {
	return source.name
}

func (source *CaldavSource) GetAuth() auth.AuthMethod {
	return source.auth
}

func (source *CaldavSource) GetSettings() primitives.SourceSettings {
	return source.settings
}

func NewCaldavSource(name string, url *types.Url, auth auth.AuthMethod) *CaldavSource {
	return &CaldavSource{
		id:   types.EmptyId(), // Placeholder until the database assigns an ID
		name: name,
		auth: auth,
		settings: &CaldavSourceSettings{
			Url: url,
		},
	}
}

func PackCaldavSource(id types.ID, name string, settings *CaldavSourceSettings, auth auth.AuthMethod) *CaldavSource {
	return &CaldavSource{
		id:       id,
		name:     name,
		settings: settings,
		auth:     auth,
	}
}

func (source *CaldavSource) getClient() (*caldav.Client, error) {
	if source.client == nil {
		var err error
		source.client, err = caldav.NewClient(
			source.auth,
			source.settings.Url.URL().String(),
		)

		if err != nil {
			return nil, err
		}
	}
	return source.client, nil
}

func (source *CaldavSource) GetCalendars() ([]primitives.Calendar, error) {
	client, err := source.getClient()
	if err != nil {
		return nil, err
	}

	cals, err := client.FindCalendars(context.TODO(), "")
	if err != nil {
		return nil, err
	}

	result := make([]primitives.Calendar, len(cals))
	for i, calendar := range cals {
		converted, err := source.calendarFromCaldav(calendar)
		if err != nil {
			return nil, fmt.Errorf("could not parse calendar %v: %w", calendar.Name, err)
		}

		casted := (primitives.Calendar)(converted)

		result[i] = casted
	}

	return result, nil
}

func (source *CaldavSource) GetCalendar(settings primitives.CalendarSettings) (primitives.Calendar, error) {
	caldavSettings := settings.(*CaldavCalendarSettings)

	client, err := source.getClient()
	if err != nil {
		return nil, err
	}

	cals, err := client.FindCalendars(context.TODO(), caldavSettings.Url.Path)
	if err != nil {
		return nil, err
	}

	if len(cals) != 1 {
		return nil, fmt.Errorf("expected exactly one calendar, got %v", len(cals))
	}

	convertedCal, err := source.calendarFromCaldav(cals[0])
	if err != nil {
		return nil, fmt.Errorf("could not convert event %v: %w", cals[0].Name, err)
	}

	castedCal := (primitives.Calendar)(convertedCal)

	return castedCal, nil
}

// TODO: Add, Edit, and Delete are not supported by upstream yet

func (source *CaldavSource) AddCalendar(name string, color *types.Color) (primitives.Calendar, error) {
	//caldavCal := calendar.(*CaldavCalendar)

	//client, err := source.getClient()
	//if err != nil {
	//	return err
	//}

	return nil, fmt.Errorf("not implemented")
}

func (source *CaldavSource) EditCalendar(calendar primitives.Calendar, name string, color *types.Color) (primitives.Calendar, error) {
	//caldavCal := calendar.(*CaldavCalendar)

	//client, err := source.getClient()
	//if err != nil {
	//	return err
	//}

	return nil, fmt.Errorf("not implemented")
}

func (source *CaldavSource) DeleteCalendar(calendar primitives.Calendar) error {
	//caldavCal := calendar.(*CaldavCalendar)

	//client, err := source.getClient()
	//if err != nil {
	//	return err
	//}

	return fmt.Errorf("not implemented")
}

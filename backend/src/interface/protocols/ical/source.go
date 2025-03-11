package ical

import (
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/interface/primitives"
	"luna-backend/types"
)

type IcalSource struct {
	id       types.ID
	name     string
	settings *IcalSourceSettings
	auth     auth.AuthMethod
}

type IcalSourceSettings struct {
	Url *types.Url `json:"url"`
}

func (settings *IcalSourceSettings) GetBytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (source *IcalSource) GetType() string {
	return types.SourceIcal
}

func (source *IcalSource) GetId() types.ID {
	return source.id
}

func (source *IcalSource) GetName() string {
	return source.name
}

func (source *IcalSource) GetAuth() auth.AuthMethod {
	return source.auth
}

func (source *IcalSource) GetSettings() primitives.SourceSettings {
	return source.settings
}

func NewIcalSource(name string, url *types.Url, auth auth.AuthMethod) *IcalSource {
	return &IcalSource{
		id:   types.EmptyId(), // Placeholder until the database assigns an ID
		name: name,
		auth: auth,
		settings: &IcalSourceSettings{
			Url: url,
		},
	}
}

func PackIcalSource(id types.ID, name string, settings *IcalSourceSettings, auth auth.AuthMethod) *IcalSource {
	return &IcalSource{
		id:       id,
		name:     name,
		settings: settings,
		auth:     auth,
	}
}

func (source *IcalSource) GetCalendars() ([]primitives.Calendar, error) {
	return nil, fmt.Errorf("not implemented")
}

func (source *IcalSource) GetCalendar(settings primitives.CalendarSettings) (primitives.Calendar, error) {
	return nil, fmt.Errorf("not implemented")
}

/* Ical source is read-only */

func (source *IcalSource) AddCalendar(name string, color *types.Color) (primitives.Calendar, error) {
	return nil, fmt.Errorf("not supported")
}

func (source *IcalSource) EditCalendar(calendar primitives.Calendar, name string, color *types.Color) (primitives.Calendar, error) {
	return nil, fmt.Errorf("not supported")
}

func (source *IcalSource) DeleteCalendar(calendar primitives.Calendar) error {
	return fmt.Errorf("not supported")
}

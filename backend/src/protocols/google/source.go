package google

import (
	"encoding/json"
	"luna-backend/constants"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
)

type GoogleSource struct {
	id       types.ID
	name     string
	settings *GoogleSourceSettings
	auth     types.AuthMethod
}

type GoogleSourceSettings struct{}

func (settings *GoogleSourceSettings) GetBytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (source *GoogleSource) GetType() string {
	return constants.SourceGoogle
}

func (source *GoogleSource) GetId() types.ID {
	return source.id
}

func (source *GoogleSource) GetName() string {
	return source.name
}

func (source *GoogleSource) GetAuth() types.AuthMethod {
	return source.auth
}

func (source *GoogleSource) GetSettings() types.SourceSettings {
	return source.settings
}

func NewGoogleSource(name string, url *types.Url, auth types.AuthMethod) *GoogleSource {
	return &GoogleSource{
		id:       types.EmptyId(), // Placeholder until the database assigns an ID
		name:     name,
		auth:     auth,
		settings: &GoogleSourceSettings{},
	}
}

func PackeGoogleSource(id types.ID, name string, settings *GoogleSourceSettings, auth types.AuthMethod) *GoogleSource {
	return &GoogleSource{
		id:       id,
		name:     name,
		settings: settings,
		auth:     auth,
	}
}

func (source *GoogleSource) GetCalendars(q types.DatabaseQueries) ([]types.Calendar, *errors.ErrorTrace) {
	return []types.Calendar{}, errors.New().Status(http.StatusNotImplemented)
}

func (source *GoogleSource) GetCalendar(settings types.CalendarSettings, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusNotImplemented)
}

func (source *GoogleSource) AddCalendar(name string, color *types.Color, _ types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusNotImplemented)
}

func (source *GoogleSource) EditCalendar(calendar types.Calendar, name string, desc string, color *types.Color, override bool, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusNotImplemented)
}

func (source *GoogleSource) DeleteCalendar(calendar types.Calendar, _ types.DatabaseQueries) *errors.ErrorTrace {
	return errors.New().Status(http.StatusNotImplemented)
}

func (source *GoogleSource) Cleanup(_ types.DatabaseQueries) *errors.ErrorTrace { return nil }

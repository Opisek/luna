package google

import (
	"encoding/json"
	"luna-backend/constants"
	"luna-backend/errors"
	"luna-backend/net"
	google "luna-backend/protocols/google/internal"
	"luna-backend/types"
	"net/http"
)

type GoogleSource struct {
	id       types.ID
	name     string
	settings *GoogleSourceSettings
	auth     types.AuthMethod
}

type GoogleSourceSettings struct {
	colors *google.Colors `json:"-"`
}

func (settings *GoogleSourceSettings) GetBytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (source *GoogleSource) getColorById(id string, isCalendar bool, q types.DatabaseQueries) (*types.Color, *types.Color, *errors.ErrorTrace) {
	if source.settings.colors == nil {
		var res google.Colors

		tr := net.FetchJson(google.ApiUrl().Subpage("colors"), "GET", source.auth, nil, "", q.GetContext(), &res)
		if tr != nil {
			return nil, nil, tr
		}

		source.settings.colors = &res
	}

	var col google.ColorDefinition
	var exists bool
	if isCalendar {
		col, exists = source.settings.colors.Calendar[id]
	} else {
		col, exists = source.settings.colors.Event[id]
	}

	if !exists {
		return nil, nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "No colors found for color id %v", id)
	}

	bgCol, err := types.ParseColor(col.Background)
	if err != nil {
		return nil, nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Could not parse background color %v for color id %v", col.Background, id)
	}

	fgCol, err := types.ParseColor(col.Foreground)
	if err != nil {
		return nil, nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Could not parse foreground color %v for color id %v", col.Background, id)
	}

	return bgCol, fgCol, nil
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

func NewGoogleSource(name string, auth types.AuthMethod) *GoogleSource {
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
	var res struct {
		Items []*google.CalendarListEntry `json:"items"`
	}

	tr := net.FetchJson(google.ApiUrl().Subpage("users", "me", "calendarList"), "GET", source.auth, nil, "", q.GetContext(), &res)
	if tr != nil {
		return nil, tr
	}

	result := make([]types.Calendar, len(res.Items))
	for i, calendar := range res.Items {
		converted, err := source.calendarFromGoogle(calendar, q)
		if err != nil {
			return nil, err.
				Append(errors.LvlBroad, "Could not get calendars")
		}

		casted := (types.Calendar)(converted)

		result[i] = casted
	}

	return result, nil
}

func (source *GoogleSource) GetCalendar(settings types.CalendarSettings, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	googleSettings := settings.(*GoogleCalendarSettings)

	var res google.CalendarListEntry

	tr := net.FetchJson(google.ApiUrl().Subpage("users", "me", "calendarList", googleSettings.GoogleId), "GET", source.auth, nil, "", q.GetContext(), &res)
	if tr != nil {
		return nil, tr
	}

	converted, err := source.calendarFromGoogle(&res, q)
	if err != nil {
		return nil, err.
			Append(errors.LvlBroad, "Could not get calendars")
	}

	casted := (types.Calendar)(converted)

	return casted, nil
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

package google

import (
	"context"
	"encoding/json"
	"luna-backend/auth"
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

	colors *google.Colors
}

type GoogleSourceSettings struct{}

func (settings *GoogleSourceSettings) GetBytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (source *GoogleSource) fetchColors(q types.DatabaseQueries) *errors.ErrorTrace {
	if source.colors == nil {
		var res google.Colors

		tr := net.FetchJson(google.ApiUrl().Subpage("colors"), "GET", source.auth, nil, "", q.GetContext(), &res)
		if tr != nil {
			return tr
		}

		source.colors = &res
	}

	return nil
}

func (source *GoogleSource) getColorById(id string, isCalendar bool, q types.DatabaseQueries) (*types.Color, *types.Color, *errors.ErrorTrace) {
	tr := source.fetchColors(q)
	if tr != nil {
		return nil, nil, nil
	}

	var col google.ColorDefinition
	var exists bool
	if isCalendar {
		col, exists = source.colors.Calendar[id]
	} else {
		col, exists = source.colors.Event[id]
	}

	if !exists {
		return nil, nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "No colors found for color id %v", id)
	}

	bgCol, err := types.ParseColor(col.Background)
	if err != nil {
		return nil, nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse background color %v for color id %v", col.Background, id)
	}

	fgCol, err := types.ParseColor(col.Foreground)
	if err != nil {
		return nil, nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse foreground color %v for color id %v", col.Background, id)
	}

	return bgCol, fgCol, nil
}

func (source *GoogleSource) getClosestColorId(col *types.Color, isCalendar bool, q types.DatabaseQueries) (string, *errors.ErrorTrace) {
	tr := source.fetchColors(q)
	if tr != nil {
		return "", tr.Append(errors.LvlDebug, "Could not map color %v to Google Calendar color id", col.String())
	}

	closestDist := ^uint(0)
	var closestCol string

	var colorMap map[string]google.ColorDefinition
	if isCalendar {
		colorMap = source.colors.Calendar
	} else {
		colorMap = source.colors.Event
	}

	// TODO: there are smarter ways to do this algorithmically, but given the low amount of colors (~6) it is fine for the time being
	for id, c := range colorMap {
		parsedGoogleColor, err := types.ParseColor(c.Background)
		if err != nil {
			return "", errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not parse background color %v for color id %v", c.Background, id).
				Append(errors.LvlDebug, "Could not map color %v to Google Calendar color id", col.String())
		}

		currentDist := col.Distance(parsedGoogleColor)
		if currentDist < closestDist {
			closestDist = currentDist
			closestCol = id
		}
	}

	return closestCol, nil
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

func (source *GoogleSource) CanAddCalendars() bool {
	return true
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

func (source *GoogleSource) AddCalendar(name string, desc string, color *types.Color, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	var tr *errors.ErrorTrace

	calendar := google.Calendar{
		Name:        name,
		Description: desc,
	}

	url := google.ApiUrl().Subpage("calendars")

	var insertedCalendar google.Calendar

	tr = net.FetchJson(url, "POST", source.auth, &calendar, "application/json", q.GetContext(), &insertedCalendar)
	if tr != nil {
		return nil, tr.
			AltStr(errors.LvlBroad, "Could not post calendar").
			Append(errors.LvlDebug, "Could not add calendar to source %v", source.GetId()).
			AltStr(errors.LvlBroad, "Could not add calendar")
	}

	var foreground string
	if color.IsDark() {
		foreground = "#ffffff"
	} else {
		foreground = "#000000"
	}

	calendarListEntry := google.CalendarListEntry{
		Name:            name,
		Description:     desc,
		BackgroundColor: color.String(),
		ForegroundColor: foreground,
	}

	url = google.ApiUrl().Subpage("users", "me", "calendarList", insertedCalendar.Id)
	query := url.Query()
	query.Set("colorRgbFormat", "true")
	url.SetQuery(query)

	var insertedCalendarListEntry google.CalendarListEntry

	tr = net.FetchJson(url, "PATCH", source.auth, &calendarListEntry, "application/json", q.GetContext(), &insertedCalendarListEntry)
	if tr != nil {
		return nil, tr.
			AltStr(errors.LvlBroad, "Could not post calendar list entry").
			Append(errors.LvlDebug, "Could not add calendar to source %v", source.GetId()).
			AltStr(errors.LvlBroad, "Could not add calendar")
	}

	converted, tr := source.calendarFromGoogle(&insertedCalendarListEntry, q)
	if tr != nil {
		return nil, tr.
			AltStr(errors.LvlBroad, "Could not parse returned calendar list entry").
			Append(errors.LvlDebug, "Could not add calendar to source %v", source.GetId()).
			AltStr(errors.LvlBroad, "Could not add calendar")
	}

	casted := (types.Calendar)(converted)

	return casted, nil
}

func (source *GoogleSource) EditCalendar(calendar types.Calendar, name string, desc string, color *types.Color, override bool, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	var foreground string
	if color.IsDark() {
		foreground = "#ffffff"
	} else {
		foreground = "#000000"
	}

	calendarListEntry := google.CalendarListEntry{
		Name:            name,
		Description:     desc,
		BackgroundColor: color.String(),
		ForegroundColor: foreground,
	}

	url := google.ApiUrl().Subpage("users", "me", "calendarList", calendar.GetSettings().(*GoogleCalendarSettings).GoogleId)
	query := url.Query()
	query.Set("colorRgbFormat", "true")
	url.SetQuery(query)

	var insertedCalendarListEntry google.CalendarListEntry

	tr := net.FetchJson(url, "PATCH", source.auth, &calendarListEntry, "application/json", q.GetContext(), &insertedCalendarListEntry)
	if tr != nil {
		return nil, tr.
			AltStr(errors.LvlBroad, "Could not patch calendar list entry").
			Append(errors.LvlDebug, "Could not edit calendar %v of source %v", calendar.GetId(), source.GetId()).
			AltStr(errors.LvlBroad, "Could not edit calendar")
	}

	converted, tr := source.calendarFromGoogle(&insertedCalendarListEntry, q)
	if tr != nil {
		return nil, tr.
			AltStr(errors.LvlBroad, "Could not parse returned calendar list entry").
			Append(errors.LvlDebug, "Could not edit calendar %v of source %v", calendar.GetId(), source.GetId()).
			AltStr(errors.LvlBroad, "Could not edit calendar")
	}

	casted := (types.Calendar)(converted)

	return casted, nil
}

func (source *GoogleSource) DeleteCalendar(calendar types.Calendar, q types.DatabaseQueries) *errors.ErrorTrace {
	googleSettings := calendar.GetSettings().(*GoogleCalendarSettings)

	url := google.ApiUrl().Subpage("calendars", googleSettings.GoogleId)

	_, tr := net.FetchBytes(url, "DELETE", source.auth, nil, "", "", q.GetContext())
	if tr != nil {
		return tr
	}

	return nil
}

func (source *GoogleSource) Cleanup(_ types.DatabaseQueries) *errors.ErrorTrace { return nil }

func (source *GoogleSource) SupplyContext(ctx context.Context) {
	if source.auth.GetType() == constants.AuthOauth {
		source.auth.(*auth.OauthAuth).SupplyContext(ctx)
	}
}

package caldav

import (
	"encoding/json"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/emersion/go-webdav/caldav"
)

type CaldavSource struct {
	id       types.ID
	name     string
	settings *CaldavSourceSettings
	auth     types.AuthMethod
}

type CaldavSourceSettings struct {
	Url    *types.Url     `json:"url"`
	client *caldav.Client `json:"-"`
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

func (source *CaldavSource) GetAuth() types.AuthMethod {
	return source.auth
}

func (source *CaldavSource) GetSettings() types.SourceSettings {
	return source.settings
}

func NewCaldavSource(name string, url *types.Url, auth types.AuthMethod) *CaldavSource {
	return &CaldavSource{
		id:   types.EmptyId(), // Placeholder until the database assigns an ID
		name: name,
		auth: auth,
		settings: &CaldavSourceSettings{
			Url: url,
		},
	}
}

func PackCaldavSource(id types.ID, name string, settings *CaldavSourceSettings, auth types.AuthMethod) *CaldavSource {
	return &CaldavSource{
		id:       id,
		name:     name,
		settings: settings,
		auth:     auth,
	}
}

func (source *CaldavSource) getClient() (*caldav.Client, *errors.ErrorTrace) {
	if source.settings.client == nil {
		var err error
		source.settings.client, err = caldav.NewClient(
			source.auth,
			source.settings.Url.URL().String(),
		)

		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not create CalDAV client")
		}
	}
	return source.settings.client, nil
}

func (source *CaldavSource) GetCalendars(q types.DatabaseQueries) ([]types.Calendar, *errors.ErrorTrace) {
	client, tr := source.getClient()
	if tr != nil {
		return nil, tr
	}

	cals, err := client.FindCalendars(q.GetContext(), "")
	if err != nil {
		return nil, errors.InterpretRemoteError(err, "source", "CalDAV source").
			Append(errors.LvlBroad, "Could not get calendars")
	}

	result := make([]types.Calendar, len(cals))
	for i, calendar := range cals {
		converted, err := source.calendarFromCaldav(calendar)
		if err != nil {
			return nil, err.
				Append(errors.LvlBroad, "Could not get calendars")
		}

		casted := (types.Calendar)(converted)

		result[i] = casted
	}

	return result, nil
}

func (source *CaldavSource) GetCalendar(settings types.CalendarSettings, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	caldavSettings := settings.(*CaldavCalendarSettings)

	client, tr := source.getClient()
	if tr != nil {
		return nil, tr
	}

	cals, err := client.FindCalendars(q.GetContext(), caldavSettings.Url.Path)
	if err != nil {
		return nil, errors.InterpretRemoteError(err, "source", "CalDAV source").
			Append(errors.LvlBroad, "Could not get calendar")
	}

	if len(cals) == 0 {
		return nil, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlBroad, "Calendar not found").
			AltStr(errors.LvlBroad, "Could not get calendar")
	}

	if len(cals) > 1 {
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlWordy, "expected exactly one calendar, got %v", len(cals)).
			Append(errors.LvlBroad, "Could not get calendar")
	}

	convertedCal, tr := source.calendarFromCaldav(cals[0])
	if tr != nil {
		return nil, tr.
			Append(errors.LvlBroad, "Could not get calendar")
	}

	castedCal := (types.Calendar)(convertedCal)

	return castedCal, nil
}

// TODO: Add, Edit, and Delete are not supported by upstream yet

func (source *CaldavSource) AddCalendar(name string, color *types.Color, _ types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	//caldavCal := calendar.(*CaldavCalendar)

	//client, err := source.getClient()
	//if err != nil {
	//	return err
	//}

	return nil, errors.New().Status(http.StatusNotImplemented)
}

func (source *CaldavSource) EditCalendar(calendar types.Calendar, name string, desc string, color *types.Color, override bool, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	if override {
		anyOverrides := false
		if name != "" {
			calendar.SetName(name)
			anyOverrides = true
		}
		if desc != "" {
			calendar.SetDesc(desc)
			anyOverrides = true
		}
		if color != nil && !color.IsEmpty() {
			calendar.SetColor(color)
			anyOverrides = true
		}

		if anyOverrides {
			q.SetCalendarOverrides(calendar.GetId(), name, desc, color)
			return calendar, nil
		} else {
			q.DeleteCalendarOverrides(calendar.GetId())
			return source.GetCalendar(calendar.GetSettings(), q)
		}
	} else {
		//caldavCal := calendar.(*CaldavCalendar)

		//client, err := source.getClient()
		//if err != nil {
		//	return err
		//}

		return nil, errors.New().Status(http.StatusNotImplemented).
			Append(errors.LvlWordy, "Only override is supported").
			Append(errors.LvlWordy, "CalDAV sources do not support editing calendars").
			AltStr(errors.LvlPlain, "This source does not support editing calendars")
	}
}

func (source *CaldavSource) DeleteCalendar(calendar types.Calendar, _ types.DatabaseQueries) *errors.ErrorTrace {
	//caldavCal := calendar.(*CaldavCalendar)

	//client, err := source.getClient()
	//if err != nil {
	//	return err
	//}

	return errors.New().Status(http.StatusNotImplemented)
}

func (source *CaldavSource) Cleanup(_ types.DatabaseQueries) *errors.ErrorTrace { return nil }

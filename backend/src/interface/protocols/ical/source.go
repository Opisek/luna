package ical

import (
	"encoding/json"
	"io"
	"luna-backend/auth"
	"luna-backend/errors"
	"luna-backend/files"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"net/http"

	"github.com/emersion/go-ical"
)

type IcalSource struct {
	id       types.ID
	name     string
	settings *IcalSourceSettings
	auth     auth.AuthMethod
}

type IcalSourceSettings struct {
	Location     string         `json:"location"`
	Url          *types.Url     `json:"url"`  // for Location == "remote"
	Path         *types.Path    `json:"path"` // for Location == "local"
	FileId       types.ID       `json:"file"` // for Location == "database"
	file         types.File     `json:"-"`
	icalCalendar *ical.Calendar `json:"-"`
}

func (source *IcalSource) getIcalFile(q types.DatabaseQueries) (*ical.Calendar, *errors.ErrorTrace) {
	if source.settings.icalCalendar == nil {
		content, tr := source.settings.file.GetContent(q)
		if tr != nil {
			return nil, tr.
				Append(errors.LvlWordy, "Could not get iCal file")
		}

		decoder := ical.NewDecoder(content)

		cal, err := decoder.Decode()
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlWordy, "Could not decode iCal file").
				AltStr(errors.LvlPlain, "Wrong file format. Are you sure the URL is correct?")
		}

		source.settings.icalCalendar = cal
	}
	return source.settings.icalCalendar, nil
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

func NewRemoteIcalSource(name string, url *types.Url, auth auth.AuthMethod) *IcalSource {
	return &IcalSource{
		id:   types.EmptyId(), // Placeholder until the database assigns an ID
		name: name,
		auth: auth,
		settings: &IcalSourceSettings{
			Location: "remote",
			Url:      url,
			file:     files.NewRemoteFile(url, auth), // TDOO: allow remote files and local (uploaded) files
		},
	}
}

func NewDatabaseIcalSource(name string, content io.Reader, q types.DatabaseQueries) (*IcalSource, *errors.ErrorTrace) {
	file, err := files.NewDatabaseFileFromContent(content, q)
	if err != nil {
		return nil, err
	}

	return &IcalSource{
		id:   types.EmptyId(), // Placeholder until the database assigns an ID
		name: name,
		auth: auth.NewNoAuth(),
		settings: &IcalSourceSettings{
			Location: "database",
			FileId:   file.GetId(),
			file:     file,
		},
	}, nil
}

func NewLocalIcalSource(name string, path *types.Path) *IcalSource {
	return &IcalSource{
		id:   types.EmptyId(), // Placeholder until the database assigns an ID
		name: name,
		auth: auth.NewNoAuth(),
		settings: &IcalSourceSettings{
			Location: "local",
			Path:     path,
			file:     files.NewLocalFile(path),
		},
	}
}

func PackIcalSource(id types.ID, name string, settings *IcalSourceSettings, auth auth.AuthMethod) (*IcalSource, *errors.ErrorTrace) {
	switch settings.Location {
	case "remote":
		settings.file = files.NewRemoteFile(settings.Url, auth)
	case "local":
		settings.file = files.NewLocalFile(settings.Path)
	case "database":
		settings.file = files.NewDatabaseFile(settings.FileId)
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).Append(errors.LvlWordy, "Unknown file location type: %v", settings.Location)
	}

	return &IcalSource{
		id:       id,
		name:     name,
		settings: settings,
		auth:     auth,
	}, nil
}

func (source *IcalSource) GetCalendars(q types.DatabaseQueries) ([]primitives.Calendar, *errors.ErrorTrace) {
	cal, err := source.getIcalFile(q)
	if err != nil {
		return nil, err.Append(errors.LvlBroad, "Could not get calendars")
	}

	result := make([]primitives.Calendar, 1)
	result[0], err = source.calendarFromIcal(cal)
	if err != nil {
		return nil, err.Append(errors.LvlBroad, "Could not get calendars")
	}

	return result, nil
}

func (source *IcalSource) GetCalendar(settings primitives.CalendarSettings, q types.DatabaseQueries) (primitives.Calendar, *errors.ErrorTrace) {
	cals, err := source.GetCalendars(q)
	if err != nil {
		return nil, err.Append(errors.LvlBroad, "Could not get calendar")
	}
	return cals[0], nil
}

/* Ical source is read-only */

func (source *IcalSource) AddCalendar(name string, color *types.Color, q types.DatabaseQueries) (primitives.Calendar, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusMethodNotAllowed)
}

func (source *IcalSource) EditCalendar(calendar primitives.Calendar, name string, color *types.Color, q types.DatabaseQueries) (primitives.Calendar, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusMethodNotAllowed)
}

func (source *IcalSource) DeleteCalendar(calendar primitives.Calendar, q types.DatabaseQueries) *errors.ErrorTrace {
	return errors.New().Status(http.StatusMethodNotAllowed)
}

func (source *IcalSource) Cleanup(q types.DatabaseQueries) *errors.ErrorTrace {
	return q.DeleteFilecache(source.settings.file)
}

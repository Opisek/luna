package ical

import (
	"context"
	"encoding/json"
	"io"
	"luna-backend/auth"
	"luna-backend/constants"
	"luna-backend/errors"
	"luna-backend/files"
	"luna-backend/types"
	"mime/multipart"
	"net/http"

	"github.com/emersion/go-ical"
)

type IcalSource struct {
	id       types.ID
	name     string
	settings *IcalSourceSettings
	auth     types.AuthMethod
}

type IcalSourceSettings struct {
	Location     string                `json:"location" form:"source_location" binding:"required"`
	Url          *types.Url            `json:"url,omitempty" form:"source_url" binding:"required_if=location remote"` // for Location == "remote"
	Path         *types.Path           `json:"path,omitempty" form:"source_path" binding:"required_if=location path"` // for Location == "local"
	File         *multipart.FileHeader `json:"-" form:"source_file" binding:"required_if=location database"`          // for Location == "database"
	FileId       types.ID              `json:"file,omitempty"`
	file         types.File            `json:"-"`
	icalCalendar *ical.Calendar        `json:"-"`
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
				AltStr(errors.LvlPlain, "Wrong file format")
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
	return constants.SourceIcal
}

func (source *IcalSource) GetId() types.ID {
	return source.id
}

func (source *IcalSource) GetName() string {
	return source.name
}

func (source *IcalSource) GetAuth() types.AuthMethod {
	return source.auth
}

func (source *IcalSource) GetSettings() types.SourceSettings {
	return source.settings
}

func (source *IcalSource) CanAddCalendars() bool {
	return false
}

func NewRemoteIcalSource(name string, url *types.Url, auth types.AuthMethod, user types.ID, q types.DatabaseQueries) (*IcalSource, *errors.ErrorTrace) {
	file, err := files.NewRemoteFile(url, "text/calendar", auth, user, q)
	if err != nil {
		return nil, err
	}

	return &IcalSource{
		id:   types.EmptyId(), // Placeholder until the database assigns an ID
		name: name,
		auth: auth,
		settings: &IcalSourceSettings{
			Location: "remote",
			Url:      url,
			file:     file,
		},
	}, nil
}

func NewDatabaseIcalSource(name string, fileName string, content io.Reader, user types.ID, q types.DatabaseQueries) (*IcalSource, *errors.ErrorTrace) {
	file, err := files.NewDatabaseFileFromContent(fileName, content, user, q)
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

func PackIcalSource(id types.ID, name string, settings *IcalSourceSettings, auth types.AuthMethod) (*IcalSource, *errors.ErrorTrace) {
	switch settings.Location {
	case "remote":
		settings.file = files.GetRemoteFile(settings.Url, "text/calendar", auth)
	case "local":
		settings.file = files.GetLocalFile(settings.Path)
	case "database":
		settings.file = files.GetDatabaseFile(settings.FileId)
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

func (source *IcalSource) GetCalendars(q types.DatabaseQueries) ([]types.Calendar, *errors.ErrorTrace) {
	cal, err := source.getIcalFile(q)
	if err != nil {
		return nil, err.Append(errors.LvlBroad, "Could not get calendars")
	}

	result := make([]types.Calendar, 1)
	result[0], err = source.calendarFromIcal(cal)
	if err != nil {
		return nil, err.Append(errors.LvlBroad, "Could not get calendars")
	}

	return result, nil
}

func (source *IcalSource) GetCalendar(settings types.CalendarSettings, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	cals, err := source.GetCalendars(q)
	if err != nil {
		return nil, err.Append(errors.LvlBroad, "Could not get calendar")
	}
	return cals[0], nil
}

/* Ical source is read-only */

func (source *IcalSource) AddCalendar(name string, desc string, color *types.Color, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusMethodNotAllowed)
}

func (source *IcalSource) EditCalendar(calendar types.Calendar, name string, desc string, color *types.Color, override bool, q types.DatabaseQueries) (types.Calendar, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusMethodNotAllowed)
}

func (source *IcalSource) DeleteCalendar(calendar types.Calendar, q types.DatabaseQueries) *errors.ErrorTrace {
	return errors.New().Status(http.StatusMethodNotAllowed)
}

func (source *IcalSource) Cleanup(q types.DatabaseQueries) *errors.ErrorTrace {
	sourceOwner, tr := q.GetSourceOwner(source.id)
	if tr != nil {
		return tr.Append(errors.LvlWordy, "Could not get file owner")
	}
	return q.DeleteFilecache(source.settings.file, sourceOwner)
}

func (source *IcalSource) SupplyContext(ctx context.Context) {
	if source.auth.GetType() == constants.AuthOauth {
		source.auth.(*auth.OauthAuth).SupplyContext(ctx)
	}
}

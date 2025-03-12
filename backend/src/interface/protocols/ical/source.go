package ical

import (
	"encoding/json"
	"fmt"
	"luna-backend/auth"
	"luna-backend/files"
	"luna-backend/interface/primitives"
	"luna-backend/types"

	"github.com/emersion/go-ical"
)

type IcalSource struct {
	id           types.ID
	name         string
	settings     *IcalSourceSettings
	auth         auth.AuthMethod
	file         types.File
	icalCalendar *ical.Calendar
}

type IcalSourceSettings struct {
	Url *types.Url `json:"url"`
}

func (source *IcalSource) getIcalFile(q types.FileQueries) (*ical.Calendar, error) {
	if source.icalCalendar == nil {
		content, err := source.file.GetContent(q)
		if err != nil {
			return nil, fmt.Errorf("could not get ical file: %v", err)
		}

		decoder := ical.NewDecoder(content)

		cal, err := decoder.Decode()
		if err != nil {
			return nil, fmt.Errorf("could not decode ical file: %v", err)
		}

		source.icalCalendar = cal
	}
	return source.icalCalendar, nil
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
		file: files.NewRemoteFile(url), // TDOO: allow remote files and local (uploaded) files
	}
}

func PackIcalSource(id types.ID, name string, settings *IcalSourceSettings, auth auth.AuthMethod) *IcalSource {
	return &IcalSource{
		id:       id,
		name:     name,
		settings: settings,
		auth:     auth,
		file:     files.NewRemoteFile(settings.Url), // TDOO: allow remote files and local (uploaded) files
	}
}

func (source *IcalSource) GetCalendars(q types.DatabaseQueries) ([]primitives.Calendar, error) {
	cal, err := source.getIcalFile(q)
	if err != nil {
		return nil, fmt.Errorf("could not get calendars: %v", err)
	}

	fmt.Println(cal.Name)

	return nil, fmt.Errorf("not implemented")
}

func (source *IcalSource) GetCalendar(settings primitives.CalendarSettings, q types.DatabaseQueries) (primitives.Calendar, error) {
	return nil, fmt.Errorf("not implemented")
}

/* Ical source is read-only */

func (source *IcalSource) AddCalendar(name string, color *types.Color, q types.DatabaseQueries) (primitives.Calendar, error) {
	return nil, fmt.Errorf("not supported")
}

func (source *IcalSource) EditCalendar(calendar primitives.Calendar, name string, color *types.Color, q types.DatabaseQueries) (primitives.Calendar, error) {
	return nil, fmt.Errorf("not supported")
}

func (source *IcalSource) DeleteCalendar(calendar primitives.Calendar, q types.DatabaseQueries) error {
	return fmt.Errorf("not supported")
}

func (source *IcalSource) Cleanup(q types.DatabaseQueries) error {
	return q.DeleteFilecache(source.file)
}

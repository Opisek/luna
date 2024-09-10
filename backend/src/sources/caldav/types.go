package caldav

import (
	"encoding/json"
	"luna-backend/auth"
	"luna-backend/sources"
	"luna-backend/types"

	"github.com/emersion/go-webdav/caldav"
)

type CaldavSource struct {
	id       types.ID
	name     string
	settings *CaldavSettings
	auth     auth.AuthMethod
	client   *caldav.Client
}

type CaldavSettings struct {
	Url *types.Url `json:"url"`
}

func (settings *CaldavSettings) GetBytes() []byte {
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

func (source *CaldavSource) GetSettings() sources.SourceSettings {
	return source.settings
}

func (source *CaldavSource) calendarFromCaldav(rawCalendar caldav.Calendar) (*types.Calendar, error) {
	return &types.Calendar{
		Source: source.GetId(),
		Id:     types.EmptyId(), // TODO: placeholder
		Path:   rawCalendar.Path,
		Name:   rawCalendar.Name,
		Desc:   rawCalendar.Description,
		Color:  nil,
	}, nil
}

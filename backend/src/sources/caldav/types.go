package caldav

import (
	"encoding/json"
	"luna-backend/auth"
	"luna-backend/sources"
	"luna-backend/types"

	"github.com/emersion/go-webdav/caldav"
)

type CaldavSource struct {
	id       sources.SourceId
	name     string
	settings *CaldavSettings
	auth     auth.AuthMethod
	client   *caldav.Client
}

type CaldavSettings struct {
	Url *types.Url `json:"url"`
}

func (source *CaldavSource) GetType() string {
	return sources.SourceCaldav
}

func (source *CaldavSource) GetId() sources.SourceId {
	return source.id
}

func (source *CaldavSource) GetName() string {
	return source.name
}

func (source *CaldavSource) GetAuth() auth.AuthMethod {
	return source.auth
}

func (source *CaldavSource) GetSettings() []byte {
	bytes, err := json.Marshal(source.settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (source *CaldavSource) calendarFromCaldav(rawCalendar caldav.Calendar) (*types.Calendar, error) {
	return &types.Calendar{
		Source: source.GetId().String(),
		Id:     rawCalendar.Path,
		Name:   rawCalendar.Name,
		Desc:   rawCalendar.Description,
		Color:  nil,
	}, nil
}

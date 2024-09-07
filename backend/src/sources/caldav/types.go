package caldav

import (
	"luna-backend/auth"
	"luna-backend/sources"
	"luna-backend/types"
	"net/url"

	"github.com/emersion/go-webdav/caldav"
)

type CaldavSource struct {
	id       *sources.SourceId
	settings *CaldavSettings
	client   *caldav.Client
}

type CaldavSettings struct {
	Url  *url.URL
	Auth auth.AuthMethod
}

func (source *CaldavSource) GetId() *sources.SourceId {
	return source.id
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

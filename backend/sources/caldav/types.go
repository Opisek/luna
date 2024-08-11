package caldav

import (
	"luna-backend/types"
	"net/url"

	"github.com/emersion/go-webdav/caldav"
)

type CaldavSource struct {
	settings *CaldavSettings
	client   *caldav.Client
}

type CaldavSettings struct {
	Url      *url.URL
	Username string
	Password string
}

func calendarFromCaldav(rawCalendar caldav.Calendar) (*types.Calendar, error) {
	return &types.Calendar{
		Name:  rawCalendar.Name,
		Desc:  rawCalendar.Description,
		Color: nil,
	}, nil
}

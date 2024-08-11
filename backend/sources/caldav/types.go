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
	calendarUrl, err := url.Parse(rawCalendar.Path)
	if err != nil {
		return nil, err
	}

	return &types.Calendar{
		Url:   calendarUrl.Host + calendarUrl.Path,
		Name:  rawCalendar.Name,
		Color: nil,
	}, nil
}

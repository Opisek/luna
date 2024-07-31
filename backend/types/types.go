package types

import (
	"image/color"
	"net/url"

	"github.com/emersion/go-webdav/caldav"
)

type Calendar struct {
	Url   string      `json:"url"`
	Name  string      `json:"name"`
	Color *color.RGBA `json:"color"`
}

func CalendarFromCaldav(rawCalendar caldav.Calendar) (*Calendar, error) {
	calendarUrl, err := url.Parse(rawCalendar.Path)
	if err != nil {
		return nil, err
	}

	return &Calendar{
		Url:   calendarUrl.Host + calendarUrl.Path,
		Name:  rawCalendar.Name,
		Color: nil,
	}, nil
}

type Event struct {
	Name string
}

type CaldavSettings struct {
	Url      *url.URL
	Username string
	Password string
}

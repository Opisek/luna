package ical

import (
	"luna-backend/types"
	"net/http"

	ics "github.com/arran4/golang-ical"
)

func (source *IcalSource) GetCalendars() ([]*types.Calendar, error) {
	res, err := source.settings.Auth.Do(&http.Request{
		Method: "GET",
		URL:    source.settings.Url,
	})
	if err != nil {
		return nil, err
	}

	cal, err := ics.ParseCalendar(res.Body)
	if err != nil {
		return nil, err
	}

	parsedCal := source.calendarFromIcal(cal)

	return []*types.Calendar{parsedCal}, nil
}

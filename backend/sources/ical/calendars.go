package ical

import (
	"luna-backend/types"
	"net/http"

	ics "github.com/arran4/golang-ical"
)

func (source *IcalSource) GetCalendars() ([]*types.Calendar, error) {
	res, err := http.Get(source.settings.Url.String())
	if err != nil {
		return nil, err
	}

	cal, err := ics.ParseCalendar(res.Body)
	if err != nil {
		return nil, err
	}

	parsedCal := calendarFromIcal(cal)

	return []*types.Calendar{parsedCal}, nil
}

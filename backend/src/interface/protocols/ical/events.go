package ical

import (
	"luna-backend/interface/primitives"
	"time"
)

func (source *IcalSource) GetEvents(calendarId string, start time.Time, end time.Time) ([]primitives.Event, error) {
	//res, err := source.settings.Auth.Do(&http.Request{
	//	Method: "GET",
	//	URL:    source.settings.Url,
	//})
	//if err != nil {
	//	return nil, err
	//}

	//cal, err := ics.ParseCalendar(res.Body)
	//if err != nil {
	//	return nil, err
	//}

	//allEvents := cal.Events()
	//filteredEvents := util.Filter(allEvents, func(event *ics.VEvent) bool {
	//	return event.
	//})

	return nil, nil
}

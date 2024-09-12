package ical

//func (source *IcalSource) GetCalendars() ([]*calendars.Calendar, error) {
//	res, err := source.settings.Auth.Do(&http.Request{
//		Method: "GET",
//		URL:    source.settings.Url.URL(),
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	cal, err := ics.ParseCalendar(res.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	parsedCal := source.calendarFromIcal(cal)
//
//	return []*calendars.Calendar{parsedCal}, nil
//}

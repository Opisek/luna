package caldav

import (
	"luna-backend/types"
	"time"
)

func (source *CaldavSource) GetEvents(calendarId string, start time.Time, end time.Time) ([]*types.Event, error) {
	//client, err := source.getClient()
	//if err != nil {
	//	return nil, err
	//}

	//rawCalendars, err := client.FindCalendars(context.TODO(), "")
	//if err != nil {
	//	return nil, err
	//}

	//var calendar *caldav.Calendar
	//for _, rawCalendar := range rawCalendars {
	//	if rawCalendar.Path == calendarId {
	//		calendar = &rawCalendar
	//		break
	//	}
	//}

	return nil, nil
}

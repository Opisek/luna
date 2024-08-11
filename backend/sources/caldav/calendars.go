package caldav

import (
	"context"
	"luna-backend/types"
)

func (source *CaldavSource) GetCalendars() ([]*types.Calendar, error) {
	client, err := source.getClient()
	if err != nil {
		return nil, err
	}

	rawCalendars, err := client.FindCalendars(context.TODO(), "")
	if err != nil {
		return nil, err
	}

	calendars := make([]*types.Calendar, len(rawCalendars))
	for i, rawCalendar := range rawCalendars {
		calendars[i], err = source.calendarFromCaldav(rawCalendar)
		if err != nil {
			return nil, err
		}
	}

	return calendars, nil
}

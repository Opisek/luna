package caldav

import (
	"context"
	"luna-backend/types"
	"net/http"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
)

func GetCalendars(caldavSettings *types.CaldavSettings) ([]*types.Calendar, error) {
	client, err := caldav.NewClient(
		webdav.HTTPClientWithBasicAuth(
			http.DefaultClient,
			caldavSettings.Username,
			caldavSettings.Password,
		),
		caldavSettings.Url.String(),
	)
	if err != nil {
		return nil, err
	}

	rawCalendars, err := client.FindCalendars(context.TODO(), "")
	if err != nil {
		return nil, err
	}

	calendars := make([]*types.Calendar, len(rawCalendars))
	for i, rawCalendar := range rawCalendars {
		calendars[i], err = types.CalendarFromCaldav(rawCalendar)
		if err != nil {
			return nil, err
		}
	}

	return calendars, nil
}

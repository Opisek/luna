package caldav

import (
	"fmt"
	"luna-backend/types"

	"github.com/emersion/go-webdav/caldav"
)

func (source *CaldavSource) calendarFromCaldav(rawCalendar caldav.Calendar) (*CaldavCalendar, error) {
	url, err := types.NewUrl(rawCalendar.Path)
	if err != nil {
		return nil, fmt.Errorf("could not parse calendar URL: %w", err)
	}

	settings := &CaldavCalendarSettings{
		Url: url,
	}

	calendar := &CaldavCalendar{
		name:     rawCalendar.Name,
		desc:     rawCalendar.Description,
		source:   source.id,
		settings: settings,
		client:   source.client,
	}

	return calendar, nil
}

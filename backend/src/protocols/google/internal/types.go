package google

import (
	"luna-backend/errors"
	"net/http"
	"time"
)

// https://developers.google.com/workspace/calendar/api/v3/reference/

type Calendar struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"summary"`
	Description string `json:"description,omitempty"`
}

type CalendarListEntry struct {
	Id              string `json:"id"`
	Name            string `json:"summary"`
	Description     string `json:"description,omitempty"`
	ColorId         string `json:"colorId,omitempty"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
	ForegroundColor string `json:"foregroundColor,omitempty"`
	Primary         bool   `json:"primary,omitempty"`
}

type ColorDefinition struct {
	Background string `json:"background"`
	Foreground string `json:"foreground"`
}

type Colors struct {
	Calendar map[string]ColorDefinition `json:"calendar"`
	Event    map[string]ColorDefinition `json:"event"`
}

type TimeDefinition struct {
	Date     string `json:"date,omitempty"`
	DateTime string `json:"dateTime,omitempty"`
	TimeZone string `json:"timeZone,omitempty"`
}

func (timeDef *TimeDefinition) String() string {
	if timeDef.Date != "" {
		return timeDef.Date
	} else {
		return timeDef.DateTime
	}
}

type Event struct {
	Id                string         `json:"id,omitempty"`
	Name              string         `json:"summary"`
	Description       string         `json:"description,omitempty"`
	ColorId           string         `json:"colorId,omitempty"`
	Start             TimeDefinition `json:"start"`
	End               TimeDefinition `json:"end"`
	Recurrence        []string       `json:"recurrence,omitempty"`
	IcalUid           string         `json:"icalUid,omitempty"`
	RecurringEventId  string         `json:"recurringEventId,omitempty"`
	Status            string         `json:"status,omitempty"`
	OriginalStartTime TimeDefinition `json:"originalStartTime,omitempty"`
}

func (timeDefinition *TimeDefinition) ParseTimeDefinition() (*time.Time, *time.Location, bool, *errors.ErrorTrace) {
	allDay := timeDefinition.Date != ""

	var parsedTime time.Time
	var err error

	timezone, err := time.LoadLocation(timeDefinition.TimeZone)
	if err != nil {
		return nil, nil, false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse timezone")
	}

	if allDay {
		if parsedTime, err = time.ParseInLocation("2006-01-02", timeDefinition.Date, timezone); err != nil {
			return nil, nil, false, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not parse date %v", timeDefinition.Date).
				AltStr(errors.LvlWordy, "Could not parse date")
		}
		// TODO: timeDefinition.timeZone
	} else {
		if parsedTime, err = time.ParseInLocation(time.RFC3339, timeDefinition.DateTime, timezone); err != nil {
			return nil, nil, false, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not parse datetime %v", timeDefinition.DateTime).
				AltStr(errors.LvlWordy, "Could not parse datetime")
		}
		// TODO: timeDefinition.timeZone
	}

	return &parsedTime, timezone, allDay, nil
}

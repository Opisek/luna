package common

import (
	"fmt"
	"luna-backend/types"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

// Custom ical properties
const (
	PropColor         = "X-LUNA-COLOR"
	PropLastColorName = "X-LUNA-LAST-COLOR-NAME"
	PropTimestamp     = "X-LUNA-TIMESTAMP"
)

func UnespaceIcalString(s string) string {
	s = strings.ReplaceAll(s, "\\,", ",")
	s = strings.ReplaceAll(s, "\\:", ",")
	s = strings.ReplaceAll(s, "\\;", ";")
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\r", "\r")
	return s
}

func EscapeIcalString(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

var IcalProductId string = "-//opisek.net//Luna//EN"

func ParseIcalTime(icalTime *ical.Prop) (*time.Time, error) {
	if icalTime == nil || icalTime.Value == "" {
		return nil, fmt.Errorf("time property is nil or empty")
	}
	timestr := icalTime.Value

	var tzid string
	if timestr[len(timestr)-1] == 'Z' {
		tzid = "UTC"
		timestr = timestr[:len(timestr)-1]
	} else {
		tzidParam := icalTime.Params.Get("TZID")
		if tzidParam == "" {
			tzid = "Local"
		} else {
			tzid = tzidParam
		}
	}

	location, err := time.LoadLocation(tzid)
	if err != nil {
		return nil, fmt.Errorf("could not parse timezone location %v: %v", tzid, err)
	}

	if !strings.Contains(timestr, "T") {
		timestr = timestr + "T000000"
	}

	parsedTime, err := time.ParseInLocation("20060102T150405", timestr, location)
	if err != nil {
		return nil, fmt.Errorf("could not parse timestamp %v: %v", timestr, err)
	}

	return &parsedTime, nil
}

type IcalEventProps struct {
	Name      string
	Desc      string
	Color     *types.Color
	Uid       string
	EventDate *types.EventDate
}

func ParseIcalEvent(props *ical.Props) (*IcalEventProps, bool, error) {
	mustUpdate := false

	// Basic info
	uid := props.Get(ical.PropUID)
	summary := props.Get(ical.PropSummary)
	summaryStr := UnespaceIcalString(summary.Value)
	description := props.Get(ical.PropDescription)
	var descStr string
	if description != nil {
		descStr = UnespaceIcalString(description.Value)
	} else {
		descStr = ""
	}

	// Color
	colorProp := props.Get("COLOR")
	lunaColorProp := props.Get(PropColor)
	lunaLastColorNameProp := props.Get(PropLastColorName)

	// If a different client has changed the color, delete the custom properties and display the new color
	if (colorProp != nil && lunaLastColorNameProp != nil && colorProp.Value != lunaLastColorNameProp.Value) || (colorProp == nil && lunaLastColorNameProp != nil) {
		lunaColorProp = nil
		lunaLastColorNameProp = nil
		mustUpdate = true
	}

	// Otherwise, parse the color normally
	var color *types.Color
	var err error
	if lunaColorProp != nil {
		color, err = types.ParseColor(lunaColorProp.Value)
	}
	if lunaColorProp == nil || err != nil {
		if colorProp == nil {
			color = types.ColorEmpty
		} else {
			color = types.ColorFromName(colorProp.Value)
			if color.IsEmpty() {
				color, err = types.ParseColor(colorProp.Value)
				if err != nil {
					color = types.ColorEmpty
				}
			}
		}
	}

	// Date
	dtstart := props.Get(ical.PropDateTimeStart)
	startTime, err := ParseIcalTime(dtstart)
	if err != nil {
		return nil, false, fmt.Errorf("could not parse start time %v: %v", dtstart.Value, err)
	}

	dtend := props.Get("DTEND")
	duration := props.Get("DURATION")

	if dtend == nil && duration == nil {
		if !(startTime.Hour() == 0 && startTime.Minute() == 0 && startTime.Second() == 0) {
			return nil, false, fmt.Errorf("event has no end time or duration")
		}
	}

	// TODO: X-CO-RECURRINGID and other ways of getting RRULE
	rrule := props.Get("RRULE")
	var eventRecurrence *types.EventRecurrence
	if rrule == nil {
		eventRecurrence = types.EmptyEventRecurrence()
	} else {
		eventRecurrence, err = types.EventRecurrenceFromIcal(rrule.Value)
		if err != nil {
			return nil, false, fmt.Errorf("could not parse recurrence rule %v: %v", rrule.Value, err)
		}
	}

	var eventDate *types.EventDate
	if dtend != nil {
		endTime, err := ParseIcalTime(dtend)
		if err != nil {
			return nil, false, fmt.Errorf("could not parse end time %v: %v", dtend.Value, err)
		}

		allDay := startTime.Location() == time.Local && endTime.Location() == time.Local && startTime.Hour() == 0 && startTime.Minute() == 0 && startTime.Second() == 0 && endTime.Hour() == 0 && endTime.Minute() == 0 && endTime.Second() == 0

		eventDate = types.NewEventDateFromEndTime(startTime, endTime, allDay, eventRecurrence)
	} else if duration != nil {
		dur, err := time.ParseDuration(duration.Value)
		if err != nil {
			return nil, false, fmt.Errorf("could not parse duration %v: %v", duration.Value, err)
		}

		allDay := startTime.Location() == time.Local && startTime.Hour() == 0 && startTime.Minute() == 0 && startTime.Second() == 0 && dur%(24*time.Hour) == 0

		eventDate = types.NewEventDateFromDuration(startTime, &dur, allDay, eventRecurrence)
	} else {
		eventDate = types.NewEventDateFromSingleDay(startTime, eventRecurrence)
	}

	parsedProps := &IcalEventProps{
		Name:      summaryStr,
		Desc:      descStr,
		Color:     color,
		Uid:       uid.Value,
		EventDate: eventDate,
	}

	return parsedProps, mustUpdate, nil
}

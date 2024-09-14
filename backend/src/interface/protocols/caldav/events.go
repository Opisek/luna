package caldav

import (
	"fmt"
	"luna-backend/crypto"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

type CaldavEvent struct {
	uid       string
	name      string
	desc      string
	settings  *CaldavEventSettings
	calendar  *CaldavCalendar
	eventDate *types.EventDate
}

type CaldavEventSettings struct {
}

func parseTime(icalTime *ical.Prop) (*time.Time, error) {
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

func eventFromCaldav(calendar *CaldavCalendar, obj *caldav.CalendarObject) (*CaldavEvent, error) {
	eventIndex := -1
	for i, child := range obj.Data.Children {
		if child.Name == "VEVENT" {
			eventIndex = i
			break
		}
	}
	if eventIndex == -1 {
		return nil, fmt.Errorf("could not find VEVENT in calendar object %v", obj.Path)
	}

	// Basic info
	uid := obj.Data.Children[eventIndex].Props.Get("UID")
	summary := obj.Data.Children[eventIndex].Props.Get("SUMMARY")
	// TODO: proper string parsing
	//summaryStr := strings.Join(strings.Split(summary.Value, "\\\\"), "\n")
	summaryStr := summary.Value
	description := obj.Data.Children[eventIndex].Props.Get("DESCRIPTION")
	var descStr string
	if description != nil {
		// TODO: proper string parsing
		//descStr = strings.Join(strings.Split(description.Value, "\\\\"), "\n")
		descStr = description.Value
	} else {
		descStr = ""
	}

	// Date
	dtstart := obj.Data.Children[eventIndex].Props.Get("DTSTART")
	startTime, err := parseTime(dtstart)
	if err != nil {
		return nil, fmt.Errorf("could not parse start time %v: %v", dtstart.Value, err)
	}

	dtend := obj.Data.Children[eventIndex].Props.Get("DTEND")
	duration := obj.Data.Children[eventIndex].Props.Get("DURATION")

	if dtend == nil && duration == nil {
		return nil, fmt.Errorf("event has no end time or duration")
	}

	// TODO: X-CO-RECURRINGID and other ways of getting RRULE
	rrule := obj.Data.Children[eventIndex].Props.Get("RRULE")
	var eventRecurrence *types.EventRecurrence
	if rrule == nil {
		eventRecurrence = types.EmptyEventRecurrence()
	} else {
		eventRecurrence, err = types.EventRecurrenceFromIcal(rrule.Value)
		if err != nil {
			return nil, fmt.Errorf("could not parse recurrence rule %v: %v", rrule.Value, err)
		}
	}

	var eventDate *types.EventDate
	if dtend != nil {
		endTime, err := parseTime(dtend)
		if err != nil {
			return nil, fmt.Errorf("could not parse end time %v: %v", dtend.Value, err)
		}
		eventDate = types.NewEventDateFromEndTime(startTime, endTime, eventRecurrence)
	} else {
		dur, err := time.ParseDuration(duration.Value)
		if err != nil {
			return nil, fmt.Errorf("could not parse duration %v: %v", duration.Value, err)
		}
		eventDate = types.NewEventDateFromDuration(startTime, &dur, eventRecurrence)
	}

	return &CaldavEvent{
		uid:       uid.Value,
		name:      summaryStr,
		desc:      descStr,
		settings:  &CaldavEventSettings{},
		calendar:  calendar,
		eventDate: eventDate,
	}, nil
}

func (event *CaldavEventSettings) GetBytes() []byte {
	return []byte{}
}

func (event *CaldavEvent) GetId() types.ID {
	return crypto.DeriveID(event.calendar.GetId(), event.uid)
}

func (event *CaldavEvent) GetName() string {
	return event.name
}

func (event *CaldavEvent) GetDesc() string {
	return event.desc
}

func (event *CaldavEvent) GetCalendar() types.ID {
	return event.calendar.GetId()
}

func (event *CaldavEvent) GetSettings() primitives.EventSettings {
	return event.settings
}

func (event *CaldavEvent) GetColor() *types.Color {
	return types.ColorFromVals(50, 50, 50)
}

func (event *CaldavEvent) GetDate() *types.EventDate {
	return event.eventDate
}

package caldav

import (
	"encoding/json"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/interface/primitives"
	util "luna-backend/interface/protocols/caldav/internal"
	"luna-backend/types"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

type CaldavEvent struct {
	name      string
	desc      string
	color     *types.Color
	settings  *CaldavEventSettings
	calendar  *CaldavCalendar
	eventDate *types.EventDate
}

type CaldavEventSettings struct {
	Url      *types.Url             `json:"url"`
	Uid      string                 `json:"uid"`
	rawEvent *caldav.CalendarObject `json:"-"`
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
	mustUpdate := false

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
	summaryStr := unespaceString(summary.Value)
	description := obj.Data.Children[eventIndex].Props.Get("DESCRIPTION")
	var descStr string
	if description != nil {
		descStr = unespaceString(description.Value)
	} else {
		descStr = ""
	}

	// Color
	colorProp := obj.Data.Children[eventIndex].Props.Get("COLOR")
	lunaColorProp := obj.Data.Children[eventIndex].Props.Get(util.PropColor)
	lunaLastColorNameProp := obj.Data.Children[eventIndex].Props.Get(util.PropLastColorName)

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

		allDay := startTime.Location() == time.Local && endTime.Location() == time.Local && startTime.Hour() == 0 && startTime.Minute() == 0 && startTime.Second() == 0 && endTime.Hour() == 0 && endTime.Minute() == 0 && endTime.Second() == 0

		eventDate = types.NewEventDateFromEndTime(startTime, endTime, allDay, eventRecurrence)
	} else {
		dur, err := time.ParseDuration(duration.Value)
		if err != nil {
			return nil, fmt.Errorf("could not parse duration %v: %v", duration.Value, err)
		}

		allDay := startTime.Location() == time.Local && startTime.Hour() == 0 && startTime.Minute() == 0 && startTime.Second() == 0 && dur%(24*time.Hour) == 0

		eventDate = types.NewEventDateFromDuration(startTime, &dur, allDay, eventRecurrence)
	}

	url, err := types.NewUrl(obj.Path)
	if err != nil {
		return nil, fmt.Errorf("could not parse event URL %v: %w", obj.Path, err)
	}

	event := &CaldavEvent{
		name:  summaryStr,
		desc:  descStr,
		color: color,
		settings: &CaldavEventSettings{
			Url:      url,
			Uid:      uid.Value,
			rawEvent: obj,
		},
		calendar:  calendar,
		eventDate: eventDate,
	}

	if mustUpdate {
		calendar.EditEvent(event, summaryStr, descStr, color, eventDate)
		// TODO: we might want to catch errors and display them as notifications here
	}

	return event, nil
}

func (settings *CaldavEventSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genEventId(calendarId types.ID, uid string) types.ID {
	return crypto.DeriveID(calendarId, uid)
}

func (event *CaldavEvent) GetId() types.ID {
	return genEventId(event.calendar.GetId(), event.settings.Uid)
}

func (event *CaldavEvent) GetName() string {
	return event.name
}

func (event *CaldavEvent) GetDesc() string {
	return event.desc
}

func (event *CaldavEvent) GetCalendar() primitives.Calendar {
	return event.calendar
}

func (event *CaldavEvent) GetSettings() primitives.EventSettings {
	return event.settings
}

func (event *CaldavEvent) GetColor() *types.Color {
	if event.color == nil {
		return event.calendar.GetColor()
	} else {
		return event.color
	}
}

func (event *CaldavEvent) SetColor(color *types.Color) {
	event.color = color
}

func (event *CaldavEvent) GetDate() *types.EventDate {
	return event.eventDate
}

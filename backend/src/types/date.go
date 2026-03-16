package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
)

type EventDate struct {
	start           *time.Time
	end             *time.Time
	duration        *time.Duration
	specifyDuration bool
	allDay          bool
	timezone        string
	timezoneOffset  int
	recurrence      *EventRecurrence
}

func NewEventDateFromEndTime(start *time.Time, end *time.Time, allDay bool, recurrence *EventRecurrence) *EventDate {
	timezone, offset := start.Zone()
	if allDay {
		newStart := start.Add(time.Duration(offset) * time.Second).UTC()

		var newEnd time.Time
		if end.After(*start) {
			_, offset = end.Zone()
			newEnd = end.Add(time.Duration(offset) * time.Second).UTC()
		} else {
			newEnd = start.Add(24 * time.Hour)
		}

		start = &newStart
		end = &newEnd
	}

	return &EventDate{
		start:           start,
		end:             end,
		allDay:          allDay,
		recurrence:      recurrence,
		timezone:        timezone,
		timezoneOffset:  offset,
		specifyDuration: false,
	}
}

func NewEventDateFromDuration(start *time.Time, duration *time.Duration, allDay bool, recurrence *EventRecurrence) *EventDate {
	timezone, offset := start.Zone()
	if allDay {
		newStart := start.Add(time.Duration(offset) * time.Second).UTC()
		start = &newStart
	}

	return &EventDate{
		start:           start,
		duration:        duration,
		allDay:          allDay,
		timezone:        timezone,
		timezoneOffset:  offset,
		recurrence:      recurrence,
		specifyDuration: true,
	}
}

func NewEventDateFromSingleDay(start *time.Time, recurrence *EventRecurrence) *EventDate {
	timezone, offset := start.Zone()
	newStart := start.Add(time.Duration(offset) * time.Second).UTC()
	newEnd := newStart

	start = &newStart
	end := &newEnd

	return &EventDate{
		start:           start,
		end:             end,
		allDay:          true,
		timezone:        timezone,
		timezoneOffset:  offset,
		recurrence:      recurrence,
		specifyDuration: false,
	}
}

func (ed *EventDate) Start() *time.Time {
	return ed.start
}

func (ed *EventDate) SetStart(start *time.Time) {
	ed.start = start
}

func (ed *EventDate) End() *time.Time {
	if ed.specifyDuration {
		endTime := ed.start.Add(*ed.duration)
		return &endTime
	} else {
		return ed.end
	}
}

func (ed *EventDate) SetEnd(end *time.Time) {
	ed.end = end
	ed.specifyDuration = false
}

func (ed *EventDate) Duration() *time.Duration {
	if ed.specifyDuration {
		return ed.duration
	} else {
		duration := ed.end.Sub(*ed.start)
		return &duration
	}
}

func (ed *EventDate) SpecifyDuration() bool {
	return ed.specifyDuration
}

func (ed *EventDate) Timezone() string {
	return ed.timezone
}

func (ed *EventDate) TimezoneOffset() int {
	return ed.timezoneOffset
}

func (ed *EventDate) SetTimezone(timezone *time.Location) {
	ed.timezone = timezone.String()
	//ed.timezoneOffset = timezone.?
}

func (ed *EventDate) Recurrence() *EventRecurrence {
	if ed.recurrence == nil {
		return EmptyEventRecurrence()
	}
	return ed.recurrence
}

func (ed *EventDate) AllDay() bool {
	return ed.allDay
}

type eventDateMarshalEnd struct {
	Start      *time.Time       `json:"start"`
	End        *time.Time       `json:"end"`
	Recurrence *EventRecurrence `json:"recurrence"`
	AllDay     bool             `json:"allDay"`
}

type eventDateMarshalDuration struct {
	Start      *time.Time       `json:"start"`
	Duration   *time.Duration   `json:"duration"`
	Recurrence *EventRecurrence `json:"recurrence"`
	AllDay     bool             `json:"allDay"`
}

func (ed *EventDate) MarshalJSON() ([]byte, error) {
	if ed.specifyDuration {
		return json.Marshal(eventDateMarshalDuration{
			Start:      ed.start,
			Duration:   ed.duration,
			Recurrence: ed.recurrence,
			AllDay:     ed.allDay,
		})
	} else {
		return json.Marshal(eventDateMarshalEnd{
			Start:      ed.start,
			End:        ed.end,
			Recurrence: ed.recurrence,
			AllDay:     ed.allDay,
		})
	}
}

func (ed *EventDate) UnmarshalJSON(data []byte) error {
	var edme eventDateMarshalEnd
	if err := json.Unmarshal(data, &edme); err == nil {
		ed.start = edme.Start
		ed.end = edme.End
		ed.recurrence = edme.Recurrence
		ed.specifyDuration = false
		return nil
	}

	var edmd eventDateMarshalDuration
	if err := json.Unmarshal(data, &edmd); err == nil {
		ed.start = edmd.Start
		ed.duration = edmd.Duration
		ed.recurrence = edmd.Recurrence
		ed.specifyDuration = true
		return nil
	}

	return nil
}

func (ed *EventDate) Clone() *EventDate {
	if ed.specifyDuration {
		return NewEventDateFromDuration(ed.start, ed.duration, ed.allDay, ed.recurrence.Clone())
	} else {
		return NewEventDateFromEndTime(ed.start, ed.end, ed.allDay, ed.recurrence.Clone())
	}
}

// RFC-5545 3.3.10, 3.8.5.3
type EventRecurrence struct {
	repeats           bool
	rule              *rrule.ROption
	exceptions        []time.Time
	modifiedInstances []time.Time
	additional        []time.Time
}

func (er *EventRecurrence) Clone() *EventRecurrence {
	if !er.repeats {
		return EmptyEventRecurrence()
	}

	return &EventRecurrence{
		repeats:    er.repeats,
		rule:       er.rule,
		exceptions: er.exceptions,
	}
}

func (er *EventRecurrence) Repeats() bool {
	return er.repeats
}

func (er *EventRecurrence) Rule() *rrule.ROption {
	return er.rule
}

func (er *EventRecurrence) Except() []time.Time {
	return er.exceptions
}

func (er *EventRecurrence) Modified() []time.Time {
	return er.modifiedInstances
}

func (er *EventRecurrence) Additional() []time.Time {
	return er.additional
}

func (er *EventRecurrence) AddException(date *time.Time) {
	er.exceptions = append(er.exceptions, *date)
}

func (er *EventRecurrence) AddModifiedInstance(date *time.Time) {
	er.modifiedInstances = append(er.modifiedInstances, *date)
}

func (er *EventRecurrence) AddAdditional(date *time.Time) {
	er.repeats = true
	er.additional = append(er.additional, *date)
}

func EmptyEventRecurrence() *EventRecurrence {
	return &EventRecurrence{
		repeats: false,
	}
}

func EventRecurrenceFromIcal(ical *ical.Props) (*EventRecurrence, error) {
	roption, err := ical.RecurrenceRule()
	if err != nil {
		return nil, fmt.Errorf("could not get recurrence rule: %v", err)
	}

	var eventRecurrence *EventRecurrence
	if roption == nil {
		eventRecurrence = EmptyEventRecurrence()
	} else {
		eventRecurrence = &EventRecurrence{
			repeats: true,
			rule:    roption,
		}
	}

	for _, prop := range ical.Values("EXDATE") {
		exceptionTime, _, err := ParseIcalTime(&prop)
		if err == nil {
			eventRecurrence.AddException(exceptionTime)
		}
	}

	for _, prop := range ical.Values("RDATE") {
		additionalTime, _, err := ParseIcalTime(&prop)
		if err == nil {
			eventRecurrence.AddAdditional(additionalTime)
		}
	}

	return eventRecurrence, nil
}

func EventRecurrenceFromLines(lines []string) (*EventRecurrence, error) {
	if len(lines) == 0 {
		return EmptyEventRecurrence(), nil
	}

	roption, err := rrule.StrToROption(strings.Join(lines, "\n"))
	if err != nil {
		return nil, fmt.Errorf("could not get recurrence rule: %v", err)
	}

	if roption == nil {
		return EmptyEventRecurrence(), nil
	}

	return &EventRecurrence{
		repeats: true,
		rule:    roption,
	}, nil
}

func ParseIcalTime(icalTime *ical.Prop) (*time.Time, *time.Location, error) {
	if icalTime == nil || icalTime.Value == "" {
		return nil, nil, fmt.Errorf("time property is nil or empty")
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
		return nil, nil, fmt.Errorf("could not parse timezone location %v: %v", tzid, err)
	}

	if !strings.Contains(timestr, "T") {
		timestr = timestr + "T000000"
	}

	parsedTime, err := time.ParseInLocation("20060102T150405", timestr, location)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse timestamp %v: %v", timestr, err)
	}

	return &parsedTime, location, nil
}

func (er EventRecurrence) MarshalJSON() ([]byte, error) {
	if er.repeats {
		return []byte(fmt.Sprintf("\"%s\"", er.rule.String())), nil
	} else {
		return []byte("false"), nil
	}
}

func (er *EventRecurrence) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "false" || str == "null" || str == "" {
		er.repeats = false
		return nil
	}

	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return fmt.Errorf("invalid recurrence rule: %s", str)
	}

	roption, err := rrule.StrToROption(str[1 : len(str)-1])
	if err != nil {
		return fmt.Errorf("could not parse recurrence rule: %v", err)
	}

	er.repeats = true
	er.rule = roption
	return nil
}

package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type EventDate struct {
	start           *time.Time
	end             *time.Time
	duration        *time.Duration
	specifyDuration bool
	allDay          bool
	recurrence      *EventRecurrence
}

func NewEventDateFromEndTime(start *time.Time, end *time.Time, allDay bool, recurrence *EventRecurrence) *EventDate {
	if allDay {
		_, offset := start.Zone()
		newStart := start.Add(time.Duration(offset) * time.Second).UTC()
		start = &newStart

		_, offset = end.Zone()
		newEnd := end.Add(time.Duration(offset) * time.Second).UTC()
		end = &newEnd
	}

	return &EventDate{
		start:           start,
		end:             end,
		allDay:          allDay,
		recurrence:      recurrence,
		specifyDuration: false,
	}
}

func NewEventDateFromDuration(start *time.Time, duration *time.Duration, allDay bool, recurrence *EventRecurrence) *EventDate {
	if allDay {
		_, offset := start.Zone()
		newStart := start.Add(time.Duration(offset) * time.Second).UTC()
		start = &newStart
	}

	return &EventDate{
		start:           start,
		duration:        duration,
		allDay:          allDay,
		recurrence:      recurrence,
		specifyDuration: true,
	}
}

func (ed *EventDate) Start() *time.Time {
	return ed.start
}

func (ed *EventDate) End() *time.Time {
	if ed.specifyDuration {
		endTime := ed.start.Add(*ed.duration)
		return &endTime
	} else {
		return ed.end
	}
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

//func (ed *EventDate) Recurrence() *EventRecurrence {
//	return ed.recurrence
//}

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

// TODO: validation according to RFC-5545 3.3.10, 3.8.5.3
type EventRecurrence struct {
	repeats bool
	rule    map[string]string
}

func EmptyEventRecurrence() *EventRecurrence {
	return &EventRecurrence{
		repeats: false,
	}
}

func EventRecurrenceFromIcal(ical string) (*EventRecurrence, error) {
	if ical == "" {
		return EmptyEventRecurrence(), nil
	}

	// PROP=VAL;PROP=VAL;...
	recMap := make(map[string]string)
	parts := strings.Split(ical, ";")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			return EmptyEventRecurrence(), fmt.Errorf("could not parse key-value pair %v", part)
		}
		recMap[kv[0]] = kv[1]
	}

	return &EventRecurrence{
		repeats: true,
		rule:    recMap,
	}, nil
}

func (er *EventRecurrence) MarshalJSON() ([]byte, error) {
	if er.repeats {
		return json.Marshal(er.rule)
	} else {
		return []byte("false"), nil
	}
}

func (er *EventRecurrence) UnmarshalJSON(data []byte) error {
	if string(data) == "false" {
		er.repeats = false
		return nil
	}

	var erMap map[string]string
	if err := json.Unmarshal(data, &erMap); err != nil {
		return fmt.Errorf("could not unmarshal EventRecurrence: %v", err)
	}

	er.repeats = true
	er.rule = erMap
	return nil
}

package types

import (
	"encoding/json"
	"fmt"
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
	recurrence      *EventRecurrence
}

func NewEventDateFromEndTime(start *time.Time, end *time.Time, allDay bool, recurrence *EventRecurrence) *EventDate {
	if allDay {
		_, offset := start.Zone()
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

func NewEventDateFromSingleDay(start *time.Time, recurrence *EventRecurrence) *EventDate {
	_, offset := start.Zone()
	newStart := start.Add(time.Duration(offset) * time.Second).UTC()
	newEnd := newStart.Add(24 * time.Hour).UTC()

	start = &newStart
	end := &newEnd

	return &EventDate{
		start:           start,
		end:             end,
		allDay:          true,
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

func (ed *EventDate) Recurrence() *EventRecurrence {
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
	repeats bool
	rule    *rrule.ROption
}

func (er *EventRecurrence) Clone() *EventRecurrence {
	if !er.repeats {
		return EmptyEventRecurrence()
	}

	return &EventRecurrence{
		repeats: er.repeats,
		rule:    er.rule,
	}
}

func (er *EventRecurrence) Repeats() bool {
	return er.repeats
}

func (er *EventRecurrence) Rule() *rrule.ROption {
	return er.rule
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

	if roption == nil {
		return EmptyEventRecurrence(), nil
	}

	return &EventRecurrence{
		repeats: true,
		rule:    roption,
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

	var roption rrule.ROption
	if err := json.Unmarshal(data, &roption); err == nil {
		er.repeats = true
		er.rule = &roption
		return nil
	}

	return nil
}

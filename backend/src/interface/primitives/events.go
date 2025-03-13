package primitives

import (
	"fmt"
	"luna-backend/types"
	"time"

	"github.com/teambition/rrule-go"
)

type Event interface {
	GetId() types.ID
	GetCalendar() Calendar
	GetName() string
	GetDesc() string
	GetColor() *types.Color
	SetColor(color *types.Color)
	GetSettings() EventSettings
	GetDate() *types.EventDate
	Clone() Event
}

type EventSettings interface {
	Bytes() []byte
}

func ExpandRecurrence(event Event, start *time.Time, end *time.Time) ([]Event, error) {
	if !event.GetDate().Recurrence().Repeats() {
		return []Event{event}, nil
	}

	r, err := rrule.NewRRule(*event.GetDate().Recurrence().Rule())
	r.DTStart(*event.GetDate().Start())
	if err != nil {
		return nil, fmt.Errorf("could not create rrule: %v", err)
	}

	timeSlices := r.Between(*start, *end, true)

	events := make([]Event, len(timeSlices))
	for i, timeSlice := range timeSlices {
		newStart := timeSlice
		newEnd := newStart.Add(*event.GetDate().Duration())

		newEvent := event.Clone()
		newEvent.GetDate().SetStart(&newStart)
		newEvent.GetDate().SetEnd(&newEnd)
		events[i] = newEvent
	}

	return events, nil
}

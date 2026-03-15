package types

import (
	"luna-backend/errors"
	"net/http"
	"time"

	"github.com/teambition/rrule-go"
)

type Event interface {
	GetId() ID
	GetCalendar() Calendar

	GetName() string
	SetName(name string)
	GetDesc() string
	SetDesc(desc string)
	GetColor() *Color
	SetColor(color *Color)
	GetOverridden() bool
	SetOverridden(overridden bool)

	CanEdit() bool
	CanDelete() bool

	GetSettings() EventSettings
	GetDate() *EventDate

	Clone() Event
	UpdateRecurrenceInstance(masterStartTime *time.Time)
}

type EventSettings interface {
	Bytes() []byte
}

func ExpandRecurrence(event Event, start *time.Time, end *time.Time) ([]Event, *errors.ErrorTrace) {
	if !event.GetDate().Recurrence().Repeats() {
		return []Event{event}, nil
	}

	r, err := rrule.NewRRule(*event.GetDate().Recurrence().Rule())
	r.DTStart(*event.GetDate().Start())
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not create RRULE for %v", event.GetId()).
			Append(errors.LvlWordy, "Could not expand event recurrence for %v", event.GetName())
	}

	rset := rrule.Set{}
	rset.RRule(r)
	for _, exception := range event.GetDate().Recurrence().Except() {
		rset.ExDate(exception)
	}
	for _, exception := range event.GetDate().Recurrence().Modified() {
		rset.ExDate(exception)
	}
	setStart := rset.GetDTStart()

	timeSlices := rset.Between(*start, *end, true)

	events := make([]Event, len(timeSlices))
	for i, timeSlice := range timeSlices {
		newStart := timeSlice
		newEnd := newStart.Add(*event.GetDate().Duration())

		newEvent := event.Clone()
		newEvent.GetDate().SetStart(&newStart)
		newEvent.GetDate().SetEnd(&newEnd)
		newEvent.UpdateRecurrenceInstance(&setStart)

		events[i] = newEvent
	}

	return events, nil
}

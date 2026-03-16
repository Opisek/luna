package types

import (
	"fmt"
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

	SupplyMasterEvent(masterEvent Event)
	GetRecurrenceId() string
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

	timezone, err := time.LoadLocation(event.GetDate().timezone)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse timezone")
	}

	rset := rrule.Set{}
	rset.RRule(r)
	exceptionSet := make(map[int64]bool)
	for _, exception := range event.GetDate().Recurrence().Except() {
		rset.ExDate(exception.In(timezone))
		exceptionSet[exception.In(timezone).Unix()] = true
	}
	for _, modified := range event.GetDate().Recurrence().Modified() {
		rset.ExDate(modified.In(timezone))
		exceptionSet[modified.In(timezone).Unix()] = true
	}
	for _, additional := range event.GetDate().Recurrence().Additional() {
		rset.RDate(additional.In(timezone))
	}

	timeSlices := rset.Between(start.In(timezone), end.In(timezone), true)

	events := make([]Event, len(timeSlices))
	actualEventCount := 0
	for _, timeSlice := range timeSlices {
		_, originalOffset := timeSlice.Zone()
		timeSlice = timeSlice.In(timezone)
		_, timezoneOffset := timeSlice.Zone()
		timeSlice = timeSlice.Add(time.Duration(originalOffset-timezoneOffset) * time.Second)

		if _, exists := exceptionSet[timeSlice.Unix()]; exists {
			continue
		}

		newStart := timeSlice
		newEnd := newStart.Add(*event.GetDate().Duration())

		newEvent := event.Clone()
		newEvent.GetDate().SetStart(&newStart)
		newEvent.GetDate().SetEnd(&newEnd)
		newEvent.SupplyMasterEvent(event)

		fmt.Println(newEvent.GetDate().Start(), newEvent.GetName())
		events[actualEventCount] = newEvent
		actualEventCount += 1
	}

	return events[:actualEventCount], nil
}

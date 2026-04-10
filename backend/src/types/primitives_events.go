package types

import (
	"luna-backend/errors"
	"time"
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
	IsRecurrenceInstance() bool
	GetRecurrenceId() string
}

type EventSettings interface {
	Bytes() []byte
}

func ExpandRecurrence(event Event, start *time.Time, end *time.Time) ([]Event, *errors.ErrorTrace) {
	if !event.GetDate().Recurrence().Repeats() || event.IsRecurrenceInstance() {
		return []Event{event}, nil
	}

	rset := event.GetDate().Recurrence().EffectiveRuleSet()
	rset.DTStart(*event.GetDate().Start())

	timezone := event.GetDate().Timezone()

	timeSlices := rset.Between(start.In(timezone), end.In(timezone), true)

	exceptionSet := make(map[int64]bool)
	for _, exception := range rset.GetExDate() {
		exceptionSet[exception.In(timezone).Unix()] = true
	}

	events := make([]Event, len(timeSlices))
	actualEventCount := 0
	for _, timeSlice := range timeSlices {
		_, originalOffset := timeSlice.Zone()
		timeSlice = timeSlice.In(timezone)
		_, timezoneOffset := timeSlice.Zone()
		timeSlice = timeSlice.Add(time.Duration(originalOffset-timezoneOffset) * time.Second)

		// For some reason this is sometimes needed
		if _, exists := exceptionSet[timeSlice.Unix()]; exists {
			continue
		}

		newStart := timeSlice
		newEnd := newStart.Add(*event.GetDate().Duration())

		newEvent := event.Clone()
		newEvent.GetDate().SetStart(&newStart)
		newEvent.GetDate().SetEnd(&newEnd)
		newEvent.SupplyMasterEvent(event)

		events[actualEventCount] = newEvent
		actualEventCount += 1
	}

	return events[:actualEventCount], nil
}

package types

import (
	"luna-backend/errors"
	"time"
)

type Calendar interface {
	GetId() ID
	GetSource() Source

	GetName() string
	SetName(name string)
	GetDesc() string
	SetDesc(desc string)
	GetColor() *Color
	SetColor(color *Color)
	GetOverridden() bool
	SetOverridden(overridden bool)

	GetSettings() CalendarSettings

	GetEvents(start time.Time, end time.Time, q DatabaseQueries) ([]Event, *errors.ErrorTrace)
	GetEvent(settings EventSettings, q DatabaseQueries) (Event, *errors.ErrorTrace)
	AddEvent(name string, desc string, color *Color, date *EventDate, q DatabaseQueries) (Event, *errors.ErrorTrace)
	EditEvent(event Event, name string, desc string, color *Color, date *EventDate, override bool, q DatabaseQueries) (Event, *errors.ErrorTrace)
	DeleteEvent(event Event, q DatabaseQueries) *errors.ErrorTrace
}

type CalendarSettings interface {
	Bytes() []byte
}

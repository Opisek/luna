package primitives

import (
	"luna-backend/errors"
	"luna-backend/types"
	"time"
)

type Calendar interface {
	GetId() types.ID
	GetSource() Source
	GetName() string
	GetDesc() string
	GetColor() *types.Color
	SetColor(color *types.Color)
	GetSettings() CalendarSettings
	GetEvents(start time.Time, end time.Time, q types.DatabaseQueries) ([]Event, *errors.ErrorTrace)
	GetEvent(settings EventSettings, q types.DatabaseQueries) (Event, *errors.ErrorTrace)
	AddEvent(name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (Event, *errors.ErrorTrace)
	EditEvent(event Event, name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (Event, *errors.ErrorTrace)
	DeleteEvent(event Event, q types.DatabaseQueries) *errors.ErrorTrace
}

type CalendarSettings interface {
	Bytes() []byte
}

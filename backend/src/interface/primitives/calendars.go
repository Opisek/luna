package primitives

import (
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
	GetEvents(start time.Time, end time.Time) ([]Event, error)
	GetEvent(settings EventSettings) (Event, error)
	AddEvent(name string, desc string, color *types.Color, date *types.EventDate) (Event, error)
	EditEvent(event Event, name string, desc string, color *types.Color, date *types.EventDate) (Event, error)
	DeleteEvent(event Event) error
}

type CalendarSettings interface {
	Bytes() []byte
}

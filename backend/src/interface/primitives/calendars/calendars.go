package calendars

import (
	"image/color"
	"luna-backend/types"
	"time"
)

type Calendar interface {
	GetId() types.ID
	GetSource() types.ID
	GetName() string
	GetColor() color.Color
	GetSettings() CalendarSettings
	GetEvents(start time.Time, end time.Time) ([]*types.Event, error)
}

type CalendarSettings interface {
	GetBytes() []byte
}

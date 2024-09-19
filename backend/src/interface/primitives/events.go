package primitives

import (
	"luna-backend/types"
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
}

type EventSettings interface {
	Bytes() []byte
}

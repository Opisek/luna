package primitives

import (
	"luna-backend/types"
	"time"
)

type Event interface {
	GetId() types.ID
	GetCalendar() types.ID
	GetName() string
	GetDesc() string
	GetColor() *types.Color
	GetSettings() EventSettings
	GetStart() time.Time
	GetEnd() time.Time
}

type EventSettings interface {
	GetBytes() []byte
}

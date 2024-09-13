package primitives

import "luna-backend/types"

type Event interface {
	GetId() types.ID
	GetCalendar() types.ID
	GetName() string
	GetDesc() string
	GetColor() *types.Color
	GetSettings() EventSettings
}

type EventSettings interface {
	GetBytes() []byte
}

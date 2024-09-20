package primitives

import (
	"luna-backend/auth"
	"luna-backend/types"
)

type Source interface {
	GetType() string
	GetId() types.ID
	GetName() string
	GetAuth() auth.AuthMethod
	GetSettings() SourceSettings
	GetCalendars() ([]Calendar, error)
	GetCalendar(settings CalendarSettings) (Calendar, error)
}

type SourceSettings interface {
	GetBytes() []byte
}

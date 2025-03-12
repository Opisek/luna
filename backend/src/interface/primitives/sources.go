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
	GetCalendars(q types.DatabaseQueries) ([]Calendar, error)
	GetCalendar(settings CalendarSettings, q types.DatabaseQueries) (Calendar, error)
	AddCalendar(name string, color *types.Color, q types.DatabaseQueries) (Calendar, error)
	EditCalendar(calendar Calendar, name string, color *types.Color, q types.DatabaseQueries) (Calendar, error)
	DeleteCalendar(calendar Calendar, q types.DatabaseQueries) error
	Cleanup(q types.DatabaseQueries) error
}

type SourceSettings interface {
	GetBytes() []byte
}

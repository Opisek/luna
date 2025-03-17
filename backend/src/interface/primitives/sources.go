package primitives

import (
	"luna-backend/auth"
	"luna-backend/errors"
	"luna-backend/types"
)

type Source interface {
	GetType() string
	GetId() types.ID
	GetName() string
	GetAuth() auth.AuthMethod
	GetSettings() SourceSettings
	GetCalendars(q types.DatabaseQueries) ([]Calendar, *errors.ErrorTrace)
	GetCalendar(settings CalendarSettings, q types.DatabaseQueries) (Calendar, *errors.ErrorTrace)
	AddCalendar(name string, color *types.Color, q types.DatabaseQueries) (Calendar, *errors.ErrorTrace)
	EditCalendar(calendar Calendar, name string, color *types.Color, q types.DatabaseQueries) (Calendar, *errors.ErrorTrace)
	DeleteCalendar(calendar Calendar, q types.DatabaseQueries) *errors.ErrorTrace
	Cleanup(q types.DatabaseQueries) *errors.ErrorTrace
}

type SourceSettings interface {
	GetBytes() []byte
}

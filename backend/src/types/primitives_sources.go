package types

import (
	"luna-backend/errors"
)

type Source interface {
	GetType() string
	GetId() ID
	GetName() string
	GetAuth() AuthMethod
	GetSettings() SourceSettings
	GetCalendars(q DatabaseQueries) ([]Calendar, *errors.ErrorTrace)
	GetCalendar(settings CalendarSettings, q DatabaseQueries) (Calendar, *errors.ErrorTrace)
	AddCalendar(name string, color *Color, q DatabaseQueries) (Calendar, *errors.ErrorTrace)
	EditCalendar(calendar Calendar, name string, desc string, color *Color, override bool, q DatabaseQueries) (Calendar, *errors.ErrorTrace)
	DeleteCalendar(calendar Calendar, q DatabaseQueries) *errors.ErrorTrace
	Cleanup(q DatabaseQueries) *errors.ErrorTrace
}

type SourceSettings interface {
	GetBytes() []byte
}

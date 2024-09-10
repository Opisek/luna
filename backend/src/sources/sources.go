package sources

import (
	"luna-backend/auth"
	"luna-backend/types"
	"time"
)

type Source interface {
	GetType() string
	GetId() types.ID
	GetName() string
	GetAuth() auth.AuthMethod
	GetSettings() SourceSettings
	GetCalendars() ([]*types.Calendar, error)
	GetEvents(calendarId string, start time.Time, end time.Time) ([]*types.Event, error)
}

type SourceSettings interface {
	GetBytes() []byte
}

package sources

import (
	"luna-backend/auth"
	"luna-backend/interface/primitives/calendars"
	"luna-backend/types"
)

type Source interface {
	GetType() string
	GetId() types.ID
	GetName() string
	GetAuth() auth.AuthMethod
	GetSettings() SourceSettings
	GetCalendars() ([]calendars.Calendar, error)
}

type SourceSettings interface {
	GetBytes() []byte
}

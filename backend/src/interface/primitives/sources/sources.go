package sources

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
	GetCalendars() ([]*types.Calendar, error)
}

type SourceSettings interface {
	GetBytes() []byte
}

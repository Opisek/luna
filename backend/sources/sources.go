package sources

import (
	"luna-backend/types"
)

type Source interface {
	GetCalendars() ([]*types.Calendar, error)
	//GetEvents(calendar *types.Calendar) ([]*types.Event, error)
}

package sources

import (
	"luna-backend/types"
	"time"

	"github.com/google/uuid"
)

type SourceId uuid.UUID

type Source interface {
	GetId() *SourceId
	GetCalendars() ([]*types.Calendar, error)
	GetEvents(calendarId string, start time.Time, end time.Time) ([]*types.Event, error)
}

func NewRandomSourceId() *SourceId {
	id, _ := uuid.NewRandom()
	sourceId := SourceId(id)
	return &sourceId
}

func (id SourceId) String() string {
	uuids := uuid.UUIDs([]uuid.UUID{uuid.UUID(id)})
	strings := uuids.Strings()
	return strings[0]
}

func SourceIdFromBytes(bytes []byte) (*SourceId, error) {
	sourceId, err := uuid.FromBytes(bytes)
	if err != nil {
		return nil, err
	}
	return (*SourceId)(&sourceId), nil
}

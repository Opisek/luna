package sources

import (
	"luna-backend/auth"
	"luna-backend/types"
	"time"

	"github.com/google/uuid"
)

type SourceId uuid.UUID

const (
	SourceCaldav = "caldav"
	SourceIcal   = "ical"
)

type Source interface {
	GetType() string
	GetId() SourceId
	GetName() string
	GetAuth() auth.AuthMethod
	GetSettings() []byte
	GetCalendars() ([]*types.Calendar, error)
	GetEvents(calendarId string, start time.Time, end time.Time) ([]*types.Event, error)
}

func NewRandomSourceId() SourceId {
	id, _ := uuid.NewRandom()
	sourceId := SourceId(id)
	return sourceId
}

func (id SourceId) String() string {
	uuids := uuid.UUIDs([]uuid.UUID{uuid.UUID(id)})
	strings := uuids.Strings()
	return strings[0]
}

func SourceIdFromBytes(bytes []byte) (SourceId, error) {
	sourceId, err := uuid.FromBytes(bytes)
	if err != nil {
		return SourceId(uuid.Nil), err
	}
	return (SourceId)(sourceId), nil
}

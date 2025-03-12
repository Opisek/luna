package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ID uuid.UUID

func EmptyId() ID {
	return ID(uuid.Nil)
}

func RandomId() ID {
	id, _ := uuid.NewRandom()
	return IdFromUuid(id)
}

func UrlNamespace() ID {
	return ID(uuid.NameSpaceURL)
}

func PathNamespace() ID {
	return ID(uuid.NameSpaceOID) // no unique well-known namespace for paths, so we'll just use OID
}

func (id ID) String() string {
	uuids := uuid.UUIDs([]uuid.UUID{uuid.UUID(id)})
	strings := uuids.Strings()
	return strings[0]
}

func (id ID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id ID) Bytes() []byte {
	return []byte(id.String())
}

func IdFromBytes(bytes []byte) (ID, error) {
	id, err := uuid.FromBytes(bytes)
	if err != nil {
		return EmptyId(), err
	}
	return IdFromUuid(id), nil
}

func IdFromString(str string) (ID, error) {
	if str == "" {
		return EmptyId(), nil
	}

	id, err := uuid.Parse(str)
	if err != nil {
		return EmptyId(), err
	}

	return IdFromUuid(id), nil
}

func IdFromUuid(uu uuid.UUID) ID {
	return ID(uu)
}

func (id ID) MarshalJSON() ([]byte, error) {
	if id == EmptyId() {
		return json.Marshal(nil)
	}
	return json.Marshal(id.String())
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var rawId string
	if err := json.Unmarshal(data, &rawId); err != nil {
		return err
	}
	ID, err := IdFromString(rawId)
	if err != nil {
		return err
	}
	*id = ID
	return nil
}

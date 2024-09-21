package crypto

import (
	"fmt"
	"luna-backend/types"

	"github.com/gofrs/uuid"
	guuid "github.com/google/uuid"
)

func RandomID() (types.ID, error) {
	uid, err := guuid.NewRandom()
	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not generate random ID: %v", err)
	}

	return types.IdFromUuid(uid), nil
}

func DeriveID(baseID types.ID, salt string) types.ID {
	baseConverted := uuid.UUID(baseID.UUID())

	combined := uuid.NewV5(baseConverted, salt)

	combinedConverted := guuid.UUID(combined)

	return types.IdFromUuid(combinedConverted)
}

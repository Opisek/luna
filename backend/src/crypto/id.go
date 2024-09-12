package crypto

import (
	"luna-backend/types"

	"github.com/gofrs/uuid"
	guuid "github.com/google/uuid"
)

func DeriveID(baseID types.ID, salt string) types.ID {
	baseConverted := uuid.UUID(baseID.UUID())

	combined := uuid.NewV5(baseConverted, salt)

	combinedConverted := guuid.UUID(combined)

	return types.IdFromUuid(combinedConverted)
}

package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"luna-backend/errors"
	"net/http"
)

func GenerateRandomBytes(n int) ([]byte, *errors.ErrorTrace) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not generate %v random bytes:", n).
			AltStr(errors.LvlWordy, "Could not generate random bytes:")
	}
	return b, nil
}

func GenerateRandomNumber() (uint32, *errors.ErrorTrace) {
	bytes, err := GenerateRandomBytes(4)
	if err != nil {
		return 0, err.
			Append(errors.LvlWordy, "Could not generate random number")
	}

	num := binary.BigEndian.Uint32(bytes)

	return num, nil
}

package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"luna-backend/errors"
	"math"
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

func GenerateRandomBase64(n int) (string, *errors.ErrorTrace) {
	bytes, err := GenerateRandomBytes((int)(math.Ceil(float64(n) / 4 * 3)))
	if err != nil {
		return "", err
	}

	str := base64.StdEncoding.EncodeToString(bytes)

	return str[:n], nil
}

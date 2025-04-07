package auth

import (
	"bytes"
	"luna-backend/constants"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"golang.org/x/crypto/argon2"
)

var DefaultAlgorithm = constants.HashArgon2
var defaultSettings = map[string]int{
	"time":    1,
	"memory":  64 * 1024,
	"threads": 4,
	"keylen":  32,
}

func PasswordStillSecure(stored *types.PasswordEntry) bool {
	// If the default algorithm changed, we want to rehash
	if stored.Algorithm != DefaultAlgorithm {
		return false
	}

	// If the default settings are stronger, we want to rehash
	for key, val := range defaultSettings {
		if stored.Parameters[key] < val {
			return false
		}
	}

	// Otherwise, the password is still strong
	return true
}

func VerifyPassword(password string, stored *types.PasswordEntry) bool {
	hashedPassword, err := hashPassword(password, stored)
	if err != nil {
		return false
	}

	return bytes.Equal(hashedPassword, stored.Hash)
}

func SecurePassword(password string) (*types.PasswordEntry, *errors.ErrorTrace) {
	ran, err := crypto.GenerateRandomBytes(16)
	if err != nil {
		return nil, err.
			Append(errors.LvlWordy, "Could not generate salt").
			Append(errors.LvlWordy, "Could not secure password")
	}

	algInfo := &types.PasswordEntry{
		Salt:       ran,
		Algorithm:  DefaultAlgorithm,
		Parameters: defaultSettings,
	}

	hash, err := hashPassword(password, algInfo)

	algInfo.Hash = hash

	return algInfo, err
}

func hashPassword(password string, algInfo *types.PasswordEntry) ([]byte, *errors.ErrorTrace) {
	switch algInfo.Algorithm {
	case constants.HashPlain:
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Plain text password storing is not allowed")
	case constants.HashArgon2:
		params := &ParametersArgon2{}

		if val, ok := algInfo.Parameters["time"]; ok {
			params.Time = uint32(val)
		} else {
			return nil, errors.New().Status(http.StatusInternalServerError).
				Append(errors.LvlDebug, "Missing time parameter")
		}

		if val, ok := algInfo.Parameters["memory"]; ok {
			params.Memory = uint32(val)
		} else {
			return nil, errors.New().Status(http.StatusInternalServerError).
				Append(errors.LvlDebug, "Missing memory parameter")
		}

		if val, ok := algInfo.Parameters["threads"]; ok {
			params.Threads = uint8(val)
		} else {
			return nil, errors.New().Status(http.StatusInternalServerError).
				Append(errors.LvlDebug, "Missing threads parameter")
		}

		if val, ok := algInfo.Parameters["keylen"]; ok {
			params.KeyLen = uint32(val)
		} else {
			return nil, errors.New().Status(http.StatusInternalServerError).
				Append(errors.LvlDebug, "Missing keylen parameter")
		}

		return argon2.IDKey([]byte(password), algInfo.Salt, params.Time, params.Memory, params.Threads, params.KeyLen), nil
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Unknown algorithm")
	}
}

type ParametersArgon2 struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

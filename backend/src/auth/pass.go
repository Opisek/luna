package auth

import (
	"errors"
	"luna-backend/types"
)

var DefaultAlgorithm = types.HashPlain

func VerifyPassword(password string, hash string, algorithm string) bool {
	hashedPassword, err := hashPassword(password, algorithm)
	if err != nil {
		return false
	}

	return hashedPassword == hash
}

func SecurePassword(password string) (string, string, error) {
	hash, err := hashPassword(password, DefaultAlgorithm)
	return hash, DefaultAlgorithm, err
}

func hashPassword(password string, algorithm string) (string, error) {
	switch algorithm {
	case types.HashPlain:
		return plain(password)
	default:
		return "", errors.New("unknown algorithm")
	}
}

// TODO: remove/disable in production code
func plain(password string) (string, error) {
	return password, nil
}

package crypto

import (
	"crypto/sha256"
)

func GetSha256Hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	digest := hash.Sum(nil)
	return digest
}

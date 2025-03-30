package crypto

import (
	"crypto/sha256"
)

func GetSha256Hash(data string) []byte {
	hash := sha256.New()
	hash.Write([]byte(data))
	digest := hash.Sum(nil)
	return digest
}

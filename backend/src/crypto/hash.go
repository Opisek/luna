package crypto

import (
	"crypto/sha256"
)

func GetSha256Hash(data ...[]byte) []byte {
	hash := sha256.New()
	for _, d := range data {
		hash.Write(d)
	}
	digest := hash.Sum(nil)
	return digest
}

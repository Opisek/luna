package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("could not generate %v random bytes: %v", n, err)
	}
	return b, nil
}

func GenerateRandomNumber() (uint32, error) {
	bytes, err := GenerateRandomBytes(4)
	if err != nil {
		return 0, fmt.Errorf("could not generate random number: %v", err)
	}

	num := binary.BigEndian.Uint32(bytes)

	return num, nil
}

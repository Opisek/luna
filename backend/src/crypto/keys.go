package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"luna-backend/common"
	"os"

	"golang.org/x/crypto/hkdf"
)

func GenerateSymmetricKey(commonConfig *common.CommonConfig, name string) ([]byte, error) {
	secret, err := GenerateRandomBytes(64)
	if err != nil {
		return nil, fmt.Errorf("could not generate random bytes: %v", err)
	}

	encodedSecret := base64.StdEncoding.EncodeToString(secret)

	path := fmt.Sprintf("%s/%s.key", commonConfig.Env.GetKeysPath(), name)
	err = os.WriteFile(path, []byte(encodedSecret), 0660)
	if err != nil {
		return nil, fmt.Errorf("could not write key file: %v", err)
	}

	return secret, nil
}

func GetSymmetricKey(commonConfig *common.CommonConfig, name string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s.key", commonConfig.Env.GetKeysPath(), name)

	_, err := os.Stat(path)
	if err == nil {
		encodedSecret, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("could not read key file: %v", err)
		}

		secret, err := base64.StdEncoding.DecodeString(string(encodedSecret))
		if err != nil {
			return nil, fmt.Errorf("could not decode key: %v", err)
		}

		return secret, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return GenerateSymmetricKey(commonConfig, name)
	} else {
		return nil, fmt.Errorf("could not check key file: %v", err)
	}
}

func DeriveKey(secret []byte, salt []byte) ([]byte, error) {
	generator := hkdf.New(sha256.New, secret, salt, nil)
	newSecret := make([]byte, 64)
	_, err := generator.Read(newSecret)
	if err != nil {
		return nil, fmt.Errorf("could not derive key: %v", err)
	}
	return newSecret, nil
}

package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"luna-backend/common"
	"luna-backend/errors"
	"net/http"
	"os"

	"golang.org/x/crypto/hkdf"
)

func GenerateSymmetricKey(commonConfig *common.CommonConfig, name string) ([]byte, *errors.ErrorTrace) {
	secret, tr := GenerateRandomBytes(64)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not generate symmetric key %v", name).
			AltStr(errors.LvlWordy, "Could not generate symmetric key")
	}

	encodedSecret := base64.StdEncoding.EncodeToString(secret)

	fileName := "%s/%s.key"
	path := fmt.Sprintf(fileName, commonConfig.Env.GetKeysPath(), name)
	err := os.WriteFile(path, []byte(encodedSecret), 0660)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not write key file %v at %v", fileName, commonConfig.Env.GetKeysPath()).
			AltStr(errors.LvlWordy, "Could not write key file").
			Append(errors.LvlDebug, "Could not generate symmetric key %v", name).
			AltStr(errors.LvlWordy, "Could not generate symmetric key")
	}

	return secret, nil
}

func GetSymmetricKey(commonConfig *common.CommonConfig, name string) ([]byte, *errors.ErrorTrace) {
	fileName := "%s/%s.key"
	path := fmt.Sprintf(fileName, commonConfig.Env.GetKeysPath(), name)

	_, err := os.Stat(path)
	if err == nil {
		encodedSecret, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not read key file %v at %v", fileName, commonConfig.Env.GetKeysPath()).
				AltStr(errors.LvlWordy, "Could not read key file").
				Append(errors.LvlDebug, "Could not get symmetric key %v", name).
				AltStr(errors.LvlWordy, "Could not get symmetric key")
		}

		secret, err := base64.StdEncoding.DecodeString(string(encodedSecret))
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not decode key file %v at %v", fileName, commonConfig.Env.GetKeysPath()).
				AltStr(errors.LvlWordy, "Could not decode key file").
				Append(errors.LvlDebug, "Could not get symmetric key %v", name).
				AltStr(errors.LvlWordy, "Could not get symmetric key")
		}

		return secret, nil
	} else if err == os.ErrNotExist {
		return GenerateSymmetricKey(commonConfig, name)
	} else {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not access key file %v at %v", fileName, commonConfig.Env.GetKeysPath()).
			AltStr(errors.LvlWordy, "Could not access key file").
			Append(errors.LvlDebug, "Could not get symmetric key %v", name).
			AltStr(errors.LvlWordy, "Could not get symmetric key")
	}
}

func DeriveKey(secret []byte, salt []byte) ([]byte, error) {
	generator := hkdf.New(sha256.New, secret, salt, nil)
	newSecret := make([]byte, 64)
	_, err := generator.Read(newSecret)
	if err != nil {
		return nil, err
	}
	return newSecret, nil
}

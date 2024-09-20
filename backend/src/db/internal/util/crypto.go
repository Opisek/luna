package util

import (
	"encoding/base64"
	"fmt"
	"luna-backend/common"
	"luna-backend/crypto"
	"luna-backend/types"
)

func GetUserEncryptionKey(commonConfig *common.CommonConfig, userId types.ID) (string, error) {
	masterKey, err := crypto.GetSymmetricKey(commonConfig, "database")
	if err != nil {
		return "", fmt.Errorf("could not get master key: %v", err)
	}
	userKey, err := crypto.DeriveKey(masterKey, userId.Bytes())
	if err != nil {
		return "", fmt.Errorf("could not derive user key: %v", err)
	}
	encodedKey := base64.StdEncoding.EncodeToString(userKey)
	return encodedKey, nil
}

func GetUserDecryptionKey(commonConfig *common.CommonConfig, userId types.ID) (string, error) {
	return GetUserEncryptionKey(commonConfig, userId)
}

package util

import (
	"encoding/base64"
	"luna-backend/common"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
)

func GetUserEncryptionKey(commonConfig *common.CommonConfig, userId types.ID) (string, *errors.ErrorTrace) {
	masterKey, tr := crypto.GetSymmetricKey(commonConfig, "database")
	if tr != nil {
		return "", tr.
			Append(errors.LvlWordy, "Could not get master key")
	}
	userKey, err := crypto.DeriveKey(masterKey, userId.Bytes())
	if err != nil {
		return "", errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not derive user key for %v", userId).
			AltStr(errors.LvlWordy, "Could not derive user key")
	}
	encodedKey := base64.StdEncoding.EncodeToString(userKey)
	return encodedKey, nil
}

func GetUserDecryptionKey(commonConfig *common.CommonConfig, userId types.ID) (string, *errors.ErrorTrace) {
	return GetUserEncryptionKey(commonConfig, userId)
}

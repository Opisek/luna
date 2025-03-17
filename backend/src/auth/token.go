package auth

import (
	"luna-backend/common"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type JsonWebToken struct {
	UserId types.ID `json:"user_id"`
	jwt.RegisteredClaims
}

func NewToken(commonConfig *common.CommonConfig, userId types.ID) (string, *errors.ErrorTrace) {
	token := JsonWebToken{UserId: userId}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token)

	key, tr := crypto.GetSymmetricKey(commonConfig, "token")
	if tr != nil {
		return "", tr
	}

	signedToken, err := jwtToken.SignedString(key)
	if err != nil {
		return "", errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not sign token")
	}

	return signedToken, nil
}

func ParseToken(commonConfig *common.CommonConfig, tokenString string) (*JsonWebToken, *errors.ErrorTrace) {
	token := &JsonWebToken{}

	_, err := jwt.ParseWithClaims(tokenString, token, func(token *jwt.Token) (interface{}, error) {
		key, tr := crypto.GetSymmetricKey(commonConfig, "token")
		if tr != nil {
			return nil, tr.SerializeError(commonConfig.DetailLevel)
		}
		return key, nil
	})

	if err != nil {
		return nil, errors.New().Status(http.StatusUnauthorized).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse token")
	}

	return token, nil
}

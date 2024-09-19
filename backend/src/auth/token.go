package auth

import (
	"fmt"
	"luna-backend/common"
	"luna-backend/crypto"
	"luna-backend/types"

	"github.com/golang-jwt/jwt/v5"
)

type JsonWebToken struct {
	UserId types.ID `json:"user_id"`
	jwt.RegisteredClaims
}

func NewToken(commonConfig *common.CommonConfig, userId types.ID) (string, error) {
	token := JsonWebToken{UserId: userId}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token)

	key, err := crypto.GetSymmetricKey(commonConfig, "token")
	if err != nil {
		return "", fmt.Errorf("could not get token key: %v", err)
	}
	return jwtToken.SignedString(key)
}

func ParseToken(commonConfig *common.CommonConfig, tokenString string) (*JsonWebToken, error) {
	token := &JsonWebToken{}

	_, err := jwt.ParseWithClaims(tokenString, token, func(token *jwt.Token) (interface{}, error) {
		key, err := crypto.GetSymmetricKey(commonConfig, "token")
		if err != nil {
			return nil, fmt.Errorf("could not get token key: %v", err)
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %v", err)
	}

	return token, err
}

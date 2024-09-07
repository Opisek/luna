package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TODO: proper key usage and algorithms in production
var secret = []byte{'s', 'e', 'c', 'r', 'e', 't'}

type JsonWebToken struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func NewToken(user uuid.UUID) (string, error) {
	token := JsonWebToken{UserId: user}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token)

	return jwtToken.SignedString(secret)
}

func ParseToken(tokenString string) (*JsonWebToken, error) {
	token := &JsonWebToken{}

	_, err := jwt.ParseWithClaims(tokenString, token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	return token, err
}
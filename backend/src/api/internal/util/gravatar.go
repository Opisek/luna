package util

import (
	"encoding/hex"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/types"
	"strings"
)

func GetDefaultProfilePictureUrl(useGravatar bool, email string) *types.Url {
	//if useGravatar {
	//	return GetGravatarUrl(email)
	//} else {

	url, err := types.NewUrl("/img/pfps/default.png")
	if err != nil {
		panic(fmt.Errorf("failed to create default profile picture URL: %v", err))
	}

	return url
}

var DefaultGravatarUrlParams string = "d=identicon"

func GetGravatarUrl(email string) *types.Url {
	return GetGravatarUrlWithParams(email, DefaultGravatarUrlParams)
}

func GetGravatarUrlWithParams(email string, params string) *types.Url {
	// Trim email from leading and trailing whitespace
	email = strings.TrimSpace(email)

	// Convert email to lowercase
	email = strings.ToLower(email)

	// Get email hash
	hash := crypto.GetSha256Hash([]byte(email))

	// Convert the hash to a hex string
	hex := hex.EncodeToString(hash)

	// Construct gravatar url
	rawUrl := "https://www.gravatar.com/avatar/" + hex + "?" + params

	// Return as URL
	// We know this will never error unless we are hit by cosmic radiation
	url, err := types.NewUrl(rawUrl)
	if err != nil {
		panic(fmt.Errorf("failed to create gravatar URL: %v", err))
	}

	return url
}

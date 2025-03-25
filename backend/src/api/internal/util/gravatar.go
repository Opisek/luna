package util

import (
	"encoding/hex"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/types"
	"strings"
)

// I had no idea whether to put this is crypto, net, parsing or where else.
// Because I expect this to only be used during registration, I will put it in
// the API's util package. For subsequent profile picture change back to
// gravatar, the frontend should generate the URL instead.
func GetGravatarUrl(email string) *types.Url {
	// Trim email from leading and trailing whitespace
	email = strings.TrimSpace(email)

	// Convert email to lowercase
	email = strings.ToLower(email)

	// Get email hash
	hash := crypto.GetSha256Hash(email)

	// Convert the hash to a hex string
	hex := hex.EncodeToString(hash)

	// Construct gravatar url
	rawUrl := "https://www.gravatar.com/avatar/" + hex

	// Return as URL
	// We know this will never error unless we are hit by cosmic radiation
	url, err := types.NewUrl(rawUrl)
	if err != nil {
		panic(fmt.Errorf("failed to create gravatar URL: %v", err))
	}

	return url
}

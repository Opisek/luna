package net

import (
	"fmt"
	"io"
	"luna-backend/auth"
	"luna-backend/types"
	"net/http"
)

func FetchFile(url *types.Url, auth auth.AuthMethod) (io.Reader, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not fetch resource: %v", err)
	}

	res, err := auth.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not fetch resource: %v", err)
	}

	return res.Body, nil
}

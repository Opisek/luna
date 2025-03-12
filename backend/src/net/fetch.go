package net

import (
	"fmt"
	"io"
	"luna-backend/types"
	"net/http"
)

func FetchFile(url *types.Url) (io.Reader, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not fetch resource: %v", err)
	}
	return resp.Body, nil
}

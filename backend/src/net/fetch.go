package net

import (
	"context"
	"io"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
)

func FetchFile(url *types.Url, auth types.AuthMethod, accept string, ctx context.Context) (io.Reader, *errors.ErrorTrace) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not create request").
			Append(errors.LvlWordy, "Could not fetch resource from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch resource")
	}

	req.Header.Set("Accept", accept)
	req = req.WithContext(ctx)

	res, err := auth.Do(req)
	if err != nil {
		return nil, errors.InterpretRemoteError(err, "file", "remote file").
			Append(errors.LvlDebug, "Could not fulfill request").
			Append(errors.LvlWordy, "Could not fetch resource from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch resource")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New().Status(res.StatusCode).
			Append(errors.LvlPlain, res.Status).
			Append(errors.LvlWordy, "Error %v", res.StatusCode).
			Append(errors.LvlDebug, "Server returned an error code").
			Append(errors.LvlWordy, "Could not fetch resource from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch resource")
	}

	return res.Body, nil
}

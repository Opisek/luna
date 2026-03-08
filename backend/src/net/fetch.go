package net

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"net/url"
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

	res, tr := auth.Do(req)
	if tr != nil {
		return nil, errors.InterpretRemoteError(tr, "file", "remote file").
			Append(errors.LvlDebug, "Could not fulfill request").
			Append(errors.LvlWordy, "Could not fetch resource from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch resource")
	}

	if res.StatusCode != http.StatusOK {
		tr := errors.New().Status(res.StatusCode)

		body, err := io.ReadAll(res.Body)
		if err == nil {
			tr.Append(errors.LvlDebug, string(body))
		}

		return nil, tr.
			Append(errors.LvlPlain, "%v", res.Status).
			Append(errors.LvlWordy, "Error %v", res.StatusCode).
			Append(errors.LvlDebug, "Server returned an error code").
			Append(errors.LvlWordy, "Could not fetch resource from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch resource")
	}

	return res.Body, nil
}

func FetchJson(url *types.Url, httpMethod string, auth types.AuthMethod, body *url.Values, bodyType string, ctx context.Context, target any) *errors.ErrorTrace {
	var payload io.Reader = nil
	if body != nil {
		payload = bytes.NewBufferString(body.Encode())
	}

	req, err := http.NewRequest(httpMethod, url.String(), payload)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not create request").
			Append(errors.LvlWordy, "Could not fetch response from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch response")
	}

	if payload != nil && bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}

	req.Header.Set("Accept", "application/json")
	req = req.WithContext(ctx)

	res, tr := auth.Do(req)
	if tr != nil {
		return errors.InterpretRemoteError(tr, "object", "object").
			Append(errors.LvlDebug, "Could not fulfill request").
			Append(errors.LvlWordy, "Could not fetch response from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch response")
	}

	if res.StatusCode != http.StatusOK {
		tr := errors.New().Status(res.StatusCode)

		body, err := io.ReadAll(res.Body)
		if err == nil {
			tr.Append(errors.LvlDebug, string(body))
		}

		return tr.
			Append(errors.LvlPlain, "%v", res.Status).
			Append(errors.LvlWordy, "Error %v", res.StatusCode).
			Append(errors.LvlDebug, "Server returned an error code").
			Append(errors.LvlWordy, "Could not fetch response from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch response")
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Could not read body").
			Append(errors.LvlWordy, "Could not fetch response from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch response")
	}

	fmt.Println(string(data))

	err = json.Unmarshal(data, target)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not unmarshal the JSON object").
			Append(errors.LvlWordy, "Could not fetch response from %v", url).
			AltStr(errors.LvlPlain, "Could not fetch response")
	}

	return nil
}

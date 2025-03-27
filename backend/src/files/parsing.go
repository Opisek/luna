package files

import (
	"io"
	"luna-backend/errors"
	"net/http"

	"github.com/emersion/go-ical"
)

func IsValidIcalFile(content io.Reader) *errors.ErrorTrace {
	decoder := ical.NewDecoder(content)

	_, err := decoder.Decode()
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not decode iCal file").
			AltStr(errors.LvlPlain, "Wrong file format")
	}

	return nil
}

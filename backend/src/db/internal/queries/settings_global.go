package queries

import (
	"luna-backend/config"
	"luna-backend/errors"
	"net/http"
)

func (q *Queries) GetRawGlobalSettings() ([]byte, *errors.ErrorTrace) {
	var settings []byte

	err := q.Tx.QueryRow(
		q.Context,
		`
		SELECT JSONB_OBJECT_AGG(key, value)
		FROM global_settings
		`,
	).Scan(&settings)

	if err != nil {
		return []byte{}, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get global settings").
			AltStr(errors.LvlPlain, "Database error")
	}

	return settings, nil
}

func (q *Queries) GetGlobalSettings() (*config.GlobalSettings, *errors.ErrorTrace) {
	var globalSettings config.GlobalSettings

	err := q.Tx.QueryRow(
		q.Context,
		`
		SELECT JSONB_OBJECT_AGG(key, value)
		FROM global_settings
		`,
	).Scan(&globalSettings)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get global settings").
			AltStr(errors.LvlPlain, "Database error")
	}

	return &globalSettings, nil
}

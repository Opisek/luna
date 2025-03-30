package queries

import (
	"luna-backend/config"
	"luna-backend/errors"
	"net/http"

	"github.com/jackc/pgx/v5"
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
			Append(errors.LvlWordy, "Could not get global settings").
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
			Append(errors.LvlWordy, "Could not get global settings").
			AltStr(errors.LvlPlain, "Database error")
	}

	return &globalSettings, nil
}

func (q *Queries) GetRawGlobalSetting(key string) ([]byte, *errors.ErrorTrace) {
	var setting []byte

	err := q.Tx.QueryRow(
		q.Context,
		`
		SELECT value
		FROM global_settings
		WHERE key = $1;
		`,
		key,
	).Scan(&setting)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return []byte{}, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlPlain, "This setting does not exist")
	default:
		return []byte{}, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get global setting").
			AltStr(errors.LvlPlain, "Database error")
	}

	return setting, nil
}

func (q *Queries) GetGlobalSetting(key string) (config.SettingsEntry, *errors.ErrorTrace) {
	setting, tr := config.GetMatchingGlobalSettingStruct(key)
	if tr != nil {
		return nil, tr
	}

	raw, tr := q.GetRawGlobalSetting(key)
	if tr != nil {
		return nil, tr
	}

	err := setting.UnmarshalJSON(raw)
	if err != nil {
		return nil, errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Invalid setting value")
	}
	return setting, nil
}

func (q *Queries) UpdateGlobalSetting(setting config.SettingsEntry) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
			UPDATE global_settings
			SET value = $1
			WHERE key = $2;
		`,
		setting,
		setting.Key(),
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not update global setting").
			AltStr(errors.LvlPlain, "Database error")
	}

	return nil
}

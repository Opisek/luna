package queries

import (
	"fmt"
	"luna-backend/config"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) InitializeUserSettings(userId types.ID) *errors.ErrorTrace {
	settings := []config.SettingsEntry{
		&config.DebugMode{},
		&config.DisplayWeekNumbers{},
		&config.FirstDayOfWeek{},
		&config.ThemeLight{},
		&config.ThemeDark{},
		&config.FontText{},
		&config.FontTime{},
		&config.DisplayAllDayEventsFilled{},
		&config.DisplayNonAllDayEventsFilled{},
		&config.DisplaySmallCalendar{},
		&config.DynamicCalendarRows{},
		&config.DynamicSmallCalendarRows{},
		&config.DisplayRoundedCorners{},
		&config.UiScaling{},
	}

	valuesString := strings.Builder{}
	keys := make([]string, len(settings))
	for i, setting := range settings {
		setting.Default()
		keys[i] = setting.Key()

		// Technically we could just inject they keys directly into the query,
		// because we know these are sanitized, but oh well, this is more "proper",
		// fits with the rest of the codebase, and you're not at danger if you're
		// not taking any chances!
		valuesString.WriteString("($1, $")
		valuesString.WriteString(strconv.Itoa(i + 2))
		valuesString.WriteString(", $")
		valuesString.WriteString(strconv.Itoa(i + 2 + len(settings)))
		valuesString.WriteString(")")
		if i != len(settings)-1 {
			valuesString.WriteString(", ")
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO user_settings (userid, key, value)
		VALUES %s;
	`, valuesString.String())

	args := make([]any, 2*len(settings)+1)
	args[0] = userId.UUID()
	for i, setting := range settings {
		args[i+1] = setting.Key()
		args[i+1+len(settings)] = setting
	}

	_, err := q.Tx.Exec(
		q.Context,
		query,
		args...,
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not initialize user's settings").
			AltStr(errors.LvlPlain, "Database error")
	}

	return nil
}

func (q *Queries) GetRawUserSettings(userId types.ID) ([]byte, *errors.ErrorTrace) {
	var settings []byte

	err := q.Tx.QueryRow(
		q.Context,
		`
		SELECT JSONB_OBJECT_AGG(key, value)
		FROM user_settings
		GROUP BY userid
		HAVING userid = $1;
		`,
		userId.UUID(),
	).Scan(&settings)

	if err != nil {
		return []byte{}, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get user's settings").
			AltStr(errors.LvlPlain, "Database error")
	}

	return settings, nil
}

func (q *Queries) GetRawUserSetting(userId types.ID, key string) ([]byte, *errors.ErrorTrace) {
	var setting []byte

	err := q.Tx.QueryRow(
		q.Context,
		`
		SELECT value
		FROM user_settings
		WHERE userid = $1 AND key = $2;
		`,
		userId.UUID(),
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
			Append(errors.LvlWordy, "Could not get user's setting").
			AltStr(errors.LvlPlain, "Database error")
	}

	return setting, nil
}

func (q *Queries) GetUserSetting(userId types.ID, key string) (config.SettingsEntry, *errors.ErrorTrace) {
	setting, tr := config.GetMatchingUserSettingStruct(key)
	if tr != nil {
		return nil, tr
	}

	raw, tr := q.GetRawUserSetting(userId, key)
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

func (q *Queries) UpdateUserSetting(userId types.ID, setting config.SettingsEntry) *errors.ErrorTrace {
	_, err := q.Tx.Exec(
		q.Context,
		`
		UPDATE user_settings
		SET value = $1
		WHERE userid = $2 AND key = $3;
		`,
		setting,
		userId.UUID(),
		setting.Key(),
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not update user's setting").
			AltStr(errors.LvlPlain, "Database error")
	}

	return nil
}

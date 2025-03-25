package queries

import (
	"fmt"
	"luna-backend/errors"
	"luna-backend/types"
	"luna-backend/types/settings"
	"net/http"
	"strconv"
	"strings"
)

func (q *Queries) InitializeUserSettings(userId types.ID) *errors.ErrorTrace {
	settings := []settings.SettingsEntry{
		&settings.DebugMode{},
		&settings.DisplayWeekNumbers{},
		&settings.FirstDayOfWeek{},
		&settings.ThemeLight{},
		&settings.ThemeDark{},
		&settings.FontText{},
		&settings.FontTime{},
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
			Append(errors.LvlDebug, "Could not initialize user's settings").
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

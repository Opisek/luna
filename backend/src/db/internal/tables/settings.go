package tables

import (
	"fmt"
	"luna-backend/config"
	"luna-backend/errors"
	"net/http"
	"strconv"
	"strings"
)

func (q *Tables) InitializeUserSettingsTable() error {
	// User settings table:
	// userid key value
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS user_settings (
			userid UUID REFERENCES users(id) ON DELETE CASCADE,
			key VARCHAR(64) NOT NULL,
			value JSONB NOT NULL,
			PRIMARY KEY (userid, key)
		);
		`,
	)

	return err
}

func (q *Tables) InitializeGlobalSettingsTable() error {
	// Global settings table:
	// key value
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS global_settings (
			key VARCHAR(64) PRIMARY KEY,
			value JSONB NOT NULL
		);
		`,
	)

	return err
}

// I would prefer this be in queries/settings_global.go, but since it's part of
// migrations, it has to be part of the Tables struct.
// TODO: might refactor that later
func (q *Tables) InitializeGlobalSettings() *errors.ErrorTrace {
	settings := config.AllDefaultGlobalSettings()

	valuesString := strings.Builder{}
	keys := make([]string, len(settings))
	for i, setting := range settings {
		setting.Default()
		keys[i] = setting.Key()

		// Technically we could just inject they keys directly into the query,
		// because we know these are sanitized, but oh well, this is more "proper",
		// fits with the rest of the codebase, and you're not at danger if you're
		// not taking any chances!
		valuesString.WriteString("($")
		valuesString.WriteString(strconv.Itoa(i + 1))
		valuesString.WriteString(", $")
		valuesString.WriteString(strconv.Itoa(i + 1 + len(settings)))
		valuesString.WriteString(")")
		if i != len(settings)-1 {
			valuesString.WriteString(", ")
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO global_settings (key, value)
		VALUES %s;
	`, valuesString.String())

	args := make([]any, 2*len(settings))
	for i, setting := range settings {
		args[i] = setting.Key()
		args[i+len(settings)] = setting
	}

	_, err := q.Tx.Exec(
		q.Context,
		query,
		args...,
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not initialize global settings").
			AltStr(errors.LvlPlain, "Database error")
	}

	return nil
}

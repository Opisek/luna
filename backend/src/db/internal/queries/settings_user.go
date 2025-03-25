package queries

import (
	"luna-backend/errors"
	"luna-backend/types"
	"luna-backend/types/settings"
)

func (q *Queries) InitializeUserSettings(userId *types.ID) *errors.ErrorTrace {
	settings := []settings.SettingsEntry{
		&settings.DebugMode{},
		&settings.DisplayWeekNumbers{},
		&settings.FirstDayOfWeek{},
		&settings.ThemeLight{},
		&settings.ThemeDark{},
		&settings.FontText{},
		&settings.FontTime{},
	}

	for _, setting := range settings {
		setting.Default()
	}

	// TODO: insert into database

	return nil
}

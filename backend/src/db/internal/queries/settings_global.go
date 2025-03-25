package queries

import (
	"luna-backend/errors"
	"luna-backend/types/settings"
)

func (q *Queries) InitializeGlobalSettings() *errors.ErrorTrace {
	settings := []settings.SettingsEntry{
		&settings.RegistrationEnabled{},
		&settings.LoggingVerbosity{},
		&settings.UseCdnFonts{},
	}

	for _, setting := range settings {
		setting.Default()
	}

	// TODO: insert into database

	return nil
}

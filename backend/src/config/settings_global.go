package config

import (
	"fmt"
	"luna-backend/common"
	"luna-backend/errors"
	"net/http"
)

const (
	KeyRegistrationEnabled = "registration_enabled"
	KeyLoggingVerbosity    = "logging_verbosity"
	KeyUseCdnFonts         = "use_cdn_fonts"
)

func GetMatchingGlobalSettingStruct(key string) (SettingsEntry, *errors.ErrorTrace) {
	switch key {
	case KeyRegistrationEnabled:
		return &RegistrationEnabled{}, nil
	case KeyLoggingVerbosity:
		return &LoggingVerbosity{}, nil
	case KeyUseCdnFonts:
		return &UseCdnFonts{}, nil
	default:
		return nil, errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Invalid setting key").
			AltStr(errors.LvlPlain, "Invalid setting name")
	}
}

func ParseGlobalSetting(key string, data []byte) (SettingsEntry, *errors.ErrorTrace) {
	entry, tr := GetMatchingGlobalSettingStruct(key)
	if tr != nil {
		return nil, tr
	}

	err := entry.UnmarshalJSON(data)
	if err != nil {
		return nil, errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Invalid setting value")
	}

	return entry, nil
}

func DefaultGlobalSetting(key string) (SettingsEntry, *errors.ErrorTrace) {
	entry, tr := GetMatchingGlobalSettingStruct(key)
	if tr != nil {
		return nil, tr
	}

	entry.Default()
	return entry, nil
}

// Whether any user can register without an explicit invitation
// Should default to false
type RegistrationEnabled struct {
	Enabled bool `json:"value"`
}

func (entry *RegistrationEnabled) Key() string {
	return KeyRegistrationEnabled
}
func (entry *RegistrationEnabled) Default() {
	entry.Enabled = false
}
func (entry *RegistrationEnabled) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *RegistrationEnabled) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// How verbose the logging should be
// Should default to 2 (LvlPlain)
type LoggingVerbosity struct {
	Verbosity int `json:"value"`
}

func (entry *LoggingVerbosity) Key() string {
	return KeyLoggingVerbosity
}
func (entry *LoggingVerbosity) Default() {
	entry.Verbosity = 2
}
func (entry *LoggingVerbosity) MarshalJSON() ([]byte, error) {
	return common.MarshalInt(entry.Verbosity), nil
}
func (entry *LoggingVerbosity) UnmarshalJSON(data []byte) error {
	verbosity, err := common.UnmarshalInt(data)
	if err != nil {
		return fmt.Errorf("could not parse verbosity level: %v", err)
	}
	if verbosity < 0 || verbosity > 3 {
		return fmt.Errorf("invalid verbosity level: %d", verbosity)
	}
	entry.Verbosity = verbosity
	return nil
}

// Whether to use Google's CDN for fonts or serve them locally
// Should default to false
type UseCdnFonts struct {
	UseCdn bool `json:"value"`
}

func (entry *UseCdnFonts) Key() string {
	return KeyUseCdnFonts
}
func (entry *UseCdnFonts) Default() {
	entry.UseCdn = false
}
func (entry *UseCdnFonts) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.UseCdn), nil
}
func (entry *UseCdnFonts) UnmarshalJSON(data []byte) (err error) {
	entry.UseCdn, err = common.UnmarshalBool(data)
	return err
}

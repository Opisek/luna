package config

import (
	"fmt"
	"luna-backend/common"
)

// Whether any user can register without an explicit invitation
// Should default to false
type RegistrationEnabled struct {
	Enabled bool
}

func (entry *RegistrationEnabled) Key() string {
	return "registration_enabled"
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
	Verbosity int
}

func (entry *LoggingVerbosity) Key() string {
	return "logging_verbosity"
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
	UseCdn bool
}

func (entry *UseCdnFonts) Key() string {
	return "use_cdn_fonts"
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

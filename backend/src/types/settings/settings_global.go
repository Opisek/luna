package settings

import (
	"fmt"
)

// Whether any user can register without an explicit invitation
// Should default to false
type RegistrationEnabled struct {
	Enabled bool
}

func (entry *RegistrationEnabled) GetKey() string {
	return "registration_enabled"
}
func (entry *RegistrationEnabled) Default() {
	entry.Enabled = false
}
func (entry *RegistrationEnabled) MarshalJSON() ([]byte, error) {
	if entry.Enabled {
		return []byte("true"), nil
	} else {
		return []byte("false"), nil
	}
}
func (entry *RegistrationEnabled) UnmarshalJSON(data []byte) error {
	if string(data) == "true" {
		entry.Enabled = true
	} else {
		entry.Enabled = false
	}
	return nil
}

// How verbose the logging should be
// Should default to 2 (LvlPlain)
type LoggingVerbosity struct {
	Verbosity int
}

func (entry *LoggingVerbosity) GetKey() string {
	return "logging_verbosity"
}
func (entry *LoggingVerbosity) Default() {
	entry.Verbosity = 2
}
func (entry *LoggingVerbosity) MarshalJSON() ([]byte, error) {
	return []byte{byte(entry.Verbosity)}, nil
}
func (entry *LoggingVerbosity) UnmarshalJSON(data []byte) error {
	entry.Verbosity = int(data[0])
	if entry.Verbosity < 0 || entry.Verbosity > 3 {
		return fmt.Errorf("invalid verbosity level: %d", entry.Verbosity)
	}
	return nil
}

// Whether to use Google's CDN for fonts or serve them locally
// Should default to false
type UseCdnFonts struct {
	UseCdn bool
}

func (entry *UseCdnFonts) GetKey() string {
	return "use_cdn_fonts"
}
func (entry *UseCdnFonts) Default() {
	entry.UseCdn = false
}
func (entry *UseCdnFonts) MarshalJSON() ([]byte, error) {
	if entry.UseCdn {
		return []byte("true"), nil
	} else {
		return []byte("false"), nil
	}
}
func (entry *UseCdnFonts) UnmarshalJSON(data []byte) error {
	if string(data) == "true" {
		entry.UseCdn = true
	} else {
		entry.UseCdn = false
	}
	return nil
}

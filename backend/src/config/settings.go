package config

type SettingsEntry interface {
	Key() string
	Default()
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// All global settings as a struct
type GlobalSettings struct {
	RegistrationEnabled RegistrationEnabled `json:"registration_enabled"`
	LoggingVerbosity    LoggingVerbosity    `json:"logging_verbosity"`
	UseCdnFonts         UseCdnFonts         `json:"use_cdn_fonts"`
}

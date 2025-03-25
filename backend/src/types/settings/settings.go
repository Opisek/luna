package settings

type SettingsEntry interface {
	GetKey() string
	Default()
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

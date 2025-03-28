package config

import (
	"fmt"
	"luna-backend/common"
	"luna-backend/errors"
	"net/http"
)

const (
	KeyDebugMode          = "debug_mode"
	KeyDisplayWeekNumbers = "display_week_numbers"
	KeyFirstDayOfWeek     = "first_day_of_week"
	KeyThemeLight         = "theme_light"
	KeyThemeDark          = "theme_dark"
	KeyFontText           = "font_text"
	KeyFontTime           = "font_time"
)

func GetMatchingUserSettingStruct(key string) (SettingsEntry, *errors.ErrorTrace) {
	switch key {
	case KeyDebugMode:
		return &DebugMode{}, nil
	case KeyDisplayWeekNumbers:
		return &DisplayWeekNumbers{}, nil
	case KeyFirstDayOfWeek:
		return &FirstDayOfWeek{}, nil
	case KeyThemeLight:
		return &ThemeLight{}, nil
	case KeyThemeDark:
		return &ThemeDark{}, nil
	case KeyFontText:
		return &FontText{}, nil
	case KeyFontTime:
		return &FontTime{}, nil
	default:
		return nil, errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Invalid setting key").
			AltStr(errors.LvlPlain, "Invalid setting name")
	}
}

func ParseUserSetting(key string, data []byte) (SettingsEntry, *errors.ErrorTrace) {
	entry, tr := GetMatchingUserSettingStruct(key)
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

func DefaultUserSetting(key string) (SettingsEntry, *errors.ErrorTrace) {
	entry, tr := GetMatchingUserSettingStruct(key)
	if tr != nil {
		return nil, tr
	}

	entry.Default()
	return entry, nil
}

// Whether the debug mode is enabled, which for example displays IDs in the UI
// Should default to false
type DebugMode struct {
	Enabled bool `json:"value"`
}

func (entry *DebugMode) Key() string {
	return "debug_mode"
}
func (entry *DebugMode) Default() {
	entry.Enabled = false
}
func (entry *DebugMode) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *DebugMode) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// Whether to display week numbers in the calendar
// Should default to false
type DisplayWeekNumbers struct {
	Enabled bool `json:"value"`
}

func (entry *DisplayWeekNumbers) Key() string {
	return "display_week_numbers"
}
func (entry *DisplayWeekNumbers) Default() {
	entry.Enabled = false
}
func (entry *DisplayWeekNumbers) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *DisplayWeekNumbers) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// The first day of the week, 0 for Sunday, 1 for Monday, etc.
// Should default to 1 (Monday)
type FirstDayOfWeek struct {
	Day int `json:"value"`
}

func (entry *FirstDayOfWeek) Key() string {
	return "first_day_of_week"
}
func (entry *FirstDayOfWeek) Default() {
	entry.Day = 1
}
func (entry *FirstDayOfWeek) MarshalJSON() ([]byte, error) {
	return common.MarshalInt(entry.Day), nil
}
func (entry *FirstDayOfWeek) UnmarshalJSON(data []byte) error {
	day, err := common.UnmarshalInt(data)
	if err != nil {
		return fmt.Errorf("could not parse day of the week: %v", err)
	}
	if day < 0 || day > 6 {
		return fmt.Errorf("invalid day of the week: %d", entry.Day)
	}
	entry.Day = day
	return nil
}

// Which theme to use for the ligth mode
// Should default to "luna-light"
type ThemeLight struct {
	Theme string `json:"value"`
}

func (entry *ThemeLight) Key() string {
	return "theme_light"
}
func (entry *ThemeLight) Default() {
	entry.Theme = "luna-light"
}
func (entry *ThemeLight) MarshalJSON() ([]byte, error) {
	return common.MarshalString(entry.Theme), nil
}
func (entry *ThemeLight) UnmarshalJSON(data []byte) (err error) {
	entry.Theme, err = common.UnmarshalString(data)
	return err
}

// Which theme to use for the dark mode
// Should default to "luna-dark"
type ThemeDark struct {
	Theme string `json:"value"`
}

func (entry *ThemeDark) Key() string {
	return "theme_dark"
}
func (entry *ThemeDark) Default() {
	entry.Theme = "luna-dark"
}
func (entry *ThemeDark) MarshalJSON() ([]byte, error) {
	return common.MarshalString(entry.Theme), nil
}
func (entry *ThemeDark) UnmarshalJSON(data []byte) (err error) {
	entry.Theme, err = common.UnmarshalString(data)
	return err
}

// Which font to use for the text
// Should default to "Atkinson Hyperlegible Next"
type FontText struct {
	Font string `json:"value"`
}

func (entry *FontText) Key() string {
	return "font_text"
}
func (entry *FontText) Default() {
	entry.Font = "Atkinson Hyperlegible Next"
}
func (entry *FontText) MarshalJSON() ([]byte, error) {
	return common.MarshalString(entry.Font), nil
}
func (entry *FontText) UnmarshalJSON(data []byte) (err error) {
	entry.Font, err = common.UnmarshalString(data)
	return err
}

// Which font to use for the time
// Should default to "Atkinson Hyperlegible Mono"
type FontTime struct {
	Font string `json:"value"`
}

func (entry *FontTime) Key() string {
	return "font_time"
}
func (entry *FontTime) Default() {
	entry.Font = "Atkinson Hyperlegible Mono"
}
func (entry *FontTime) MarshalJSON() ([]byte, error) {
	return common.MarshalString(entry.Font), nil
}
func (entry *FontTime) UnmarshalJSON(data []byte) (err error) {
	entry.Font, err = common.UnmarshalString(data)
	return err
}

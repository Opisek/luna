package settings

import (
	"fmt"
	"luna-backend/parsing"
)

// Whether the debug mode is enabled, which for example displays IDs in the UI
// Should default to false
type DebugMode struct {
	Enabled bool
}

func (entry DebugMode) Key() string {
	return "debug_mode"
}
func (entry *DebugMode) Default() {
	entry.Enabled = false
}
func (entry *DebugMode) MarshalJSON() ([]byte, error) {
	return parsing.MarshalBool(entry.Enabled), nil
}
func (entry *DebugMode) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = parsing.UnmarshalBool(data)
	return err
}

// Whether to display week numbers in the calendar
// Should default to false
type DisplayWeekNumbers struct {
	Enabled bool
}

func (entry *DisplayWeekNumbers) Key() string {
	return "display_week_numbers"
}
func (entry *DisplayWeekNumbers) Default() {
	entry.Enabled = false
}
func (entry *DisplayWeekNumbers) MarshalJSON() ([]byte, error) {
	return parsing.MarshalBool(entry.Enabled), nil
}
func (entry *DisplayWeekNumbers) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = parsing.UnmarshalBool(data)
	return err
}

// The first day of the week, 0 for Sunday, 1 for Monday, etc.
// Should default to 1 (Monday)
type FirstDayOfWeek struct {
	Day int
}

func (entry *FirstDayOfWeek) Key() string {
	return "first_day_of_week"
}
func (entry *FirstDayOfWeek) Default() {
	entry.Day = 1
}
func (entry *FirstDayOfWeek) MarshalJSON() ([]byte, error) {
	return parsing.MarshalInt(entry.Day), nil
}
func (entry *FirstDayOfWeek) UnmarshalJSON(data []byte) error {
	day, err := parsing.UnmarshalInt(data)
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
	Theme string
}

func (entry *ThemeLight) Key() string {
	return "theme_light"
}
func (entry *ThemeLight) Default() {
	entry.Theme = "luna-light"
}
func (entry *ThemeLight) MarshalJSON() ([]byte, error) {
	return parsing.MarshalString(entry.Theme), nil
}
func (entry *ThemeLight) UnmarshalJSON(data []byte) (err error) {
	entry.Theme, err = parsing.UnmarshalString(data)
	return err
}

// Which theme to use for the dark mode
// Should default to "luna-dark"
type ThemeDark struct {
	Theme string
}

func (entry *ThemeDark) Key() string {
	return "theme_dark"
}
func (entry *ThemeDark) Default() {
	entry.Theme = "luna-dark"
}
func (entry *ThemeDark) MarshalJSON() ([]byte, error) {
	return parsing.MarshalString(entry.Theme), nil
}
func (entry *ThemeDark) UnmarshalJSON(data []byte) (err error) {
	entry.Theme, err = parsing.UnmarshalString(data)
	return err
}

// Which font to use for the text
// Should default to "Atkinson Hyperlegible Next"
type FontText struct {
	Font string
}

func (entry *FontText) Key() string {
	return "font_text"
}
func (entry *FontText) Default() {
	entry.Font = "Atkinson Hyperlegible Next"
}
func (entry *FontText) MarshalJSON() ([]byte, error) {
	return parsing.MarshalString(entry.Font), nil
}
func (entry *FontText) UnmarshalJSON(data []byte) (err error) {
	entry.Font, err = parsing.UnmarshalString(data)
	return err
}

// Which font to use for the time
// Should default to "Atkinson Hyperlegible Mono"
type FontTime struct {
	Font string
}

func (entry *FontTime) Key() string {
	return "font_time"
}
func (entry *FontTime) Default() {
	entry.Font = "Atkinson Hyperlegible Mono"
}
func (entry *FontTime) MarshalJSON() ([]byte, error) {
	return parsing.MarshalString(entry.Font), nil
}
func (entry *FontTime) UnmarshalJSON(data []byte) (err error) {
	entry.Font, err = parsing.UnmarshalString(data)
	return err
}

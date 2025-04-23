package config

import (
	"fmt"
	"luna-backend/common"
	"luna-backend/errors"
	"net/http"
)

const (
	KeyDebugMode                    = "debug_mode"
	KeyDisplayWeekNumbers           = "display_week_numbers"
	KeyFirstDayOfWeek               = "first_day_of_week"
	KeyThemeLight                   = "theme_light"
	KeyThemeDark                    = "theme_dark"
	KeyFontText                     = "font_text"
	KeyFontTime                     = "font_time"
	KeyDisplayAllDayEventsFilled    = "display_all_day_events_filled"
	KeyDisplayNonAllDayEventsFilled = "display_non_all_day_events_filled"
	KeyDisplaySmallCalendar         = "display_small_calendar"
	KeyDynamicCalendarRows          = "dynamic_calendar_rows"
	KeyDynamicSmallCalendarRows     = "dynamic_small_calendar_rows"
	KeyDisplayRoundedCorners        = "display_rounded_corners"
	KeyUiScaling                    = "ui_scaling"
	KeyAnimateCalendarSwipe         = "animate_calendar_swipe"
	KeyAnimateSmallCalendarSwipe    = "animate_small_calendar_swipe"
	KeyAnimateMonthSelectionSwipe   = "animate_month_selection_swipe"
)

func AllDefaultUserSettings() []SettingsEntry {
	settings := []SettingsEntry{
		&DebugMode{},
		&DisplayWeekNumbers{},
		&FirstDayOfWeek{},
		&ThemeLight{},
		&ThemeDark{},
		&FontText{},
		&FontTime{},
		&DisplayAllDayEventsFilled{},
		&DisplayNonAllDayEventsFilled{},
		&DisplaySmallCalendar{},
		&DynamicCalendarRows{},
		&DynamicSmallCalendarRows{},
		&DisplayRoundedCorners{},
		&UiScaling{},
		&AnimateCalendarSwipe{},
		&AnimateSmallCalendarSwipe{},
		&AnimateMonthSelectionSwipe{},
	}

	for _, setting := range settings {
		setting.Default()
	}

	return settings
}

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
	case KeyDisplayAllDayEventsFilled:
		return &DisplayAllDayEventsFilled{}, nil
	case KeyDisplayNonAllDayEventsFilled:
		return &DisplayNonAllDayEventsFilled{}, nil
	case KeyDisplaySmallCalendar:
		return &DisplaySmallCalendar{}, nil
	case KeyDynamicCalendarRows:
		return &DynamicCalendarRows{}, nil
	case KeyDynamicSmallCalendarRows:
		return &DynamicSmallCalendarRows{}, nil
	case KeyDisplayRoundedCorners:
		return &DisplayRoundedCorners{}, nil
	case KeyUiScaling:
		return &UiScaling{}, nil
	case KeyAnimateCalendarSwipe:
		return &AnimateCalendarSwipe{}, nil
	case KeyAnimateSmallCalendarSwipe:
		return &AnimateSmallCalendarSwipe{}, nil
	case KeyAnimateMonthSelectionSwipe:
		return &AnimateMonthSelectionSwipe{}, nil
	default:
		return nil, errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Invalid setting key %s", key).
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

// Whether to display all day events filled-in with color
// Should default to true
type DisplayAllDayEventsFilled struct {
	Enabled bool `json:"value"`
}

func (entry *DisplayAllDayEventsFilled) Key() string {
	return KeyDisplayAllDayEventsFilled
}
func (entry *DisplayAllDayEventsFilled) Default() {
	entry.Enabled = true
}
func (entry *DisplayAllDayEventsFilled) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *DisplayAllDayEventsFilled) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// Whether to display non-all day events filled-in with color
// Should default to false
type DisplayNonAllDayEventsFilled struct {
	Enabled bool `json:"value"`
}

func (entry *DisplayNonAllDayEventsFilled) Key() string {
	return KeyDisplayNonAllDayEventsFilled
}
func (entry *DisplayNonAllDayEventsFilled) Default() {
	entry.Enabled = false
}
func (entry *DisplayNonAllDayEventsFilled) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *DisplayNonAllDayEventsFilled) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// Whether to display the small calendar
// Should default to true
type DisplaySmallCalendar struct {
	Enabled bool `json:"value"`
}

func (entry *DisplaySmallCalendar) Key() string {
	return KeyDisplaySmallCalendar
}
func (entry *DisplaySmallCalendar) Default() {
	entry.Enabled = true
}
func (entry *DisplaySmallCalendar) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *DisplaySmallCalendar) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// Whether to use dynamic calendar rows
// Should default to true
type DynamicCalendarRows struct {
	Enabled bool `json:"value"`
}

func (entry *DynamicCalendarRows) Key() string {
	return KeyDynamicCalendarRows
}
func (entry *DynamicCalendarRows) Default() {
	entry.Enabled = true
}
func (entry *DynamicCalendarRows) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *DynamicCalendarRows) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// Whether to use dynamic small calendar rows
// Should default to false
type DynamicSmallCalendarRows struct {
	Enabled bool `json:"value"`
}

func (entry *DynamicSmallCalendarRows) Key() string {
	return KeyDynamicSmallCalendarRows
}
func (entry *DynamicSmallCalendarRows) Default() {
	entry.Enabled = false
}
func (entry *DynamicSmallCalendarRows) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *DynamicSmallCalendarRows) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// Whether to display rounded corners
// Should default to true
type DisplayRoundedCorners struct {
	Enabled bool `json:"value"`
}

func (entry *DisplayRoundedCorners) Key() string {
	return KeyDisplayRoundedCorners
}
func (entry *DisplayRoundedCorners) Default() {
	entry.Enabled = true
}
func (entry *DisplayRoundedCorners) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *DisplayRoundedCorners) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// The factor by which to scale the UI
// Should default to 1.0
type UiScaling struct {
	Factor float64 `json:"value"`
}

func (entry *UiScaling) Key() string {
	return KeyUiScaling
}
func (entry *UiScaling) Default() {
	entry.Factor = 1.0
}
func (entry *UiScaling) MarshalJSON() ([]byte, error) {
	return common.MarshalFloat(entry.Factor), nil
}
func (entry *UiScaling) UnmarshalJSON(data []byte) (err error) {
	entry.Factor, err = common.UnmarshalFloat(data)
	if err != nil {
		return fmt.Errorf("could not parse scaling factor: %v", err)
	}
	if entry.Factor < 0.5 || entry.Factor > 2.0 {
		return fmt.Errorf("invalid scaling factor: %f", entry.Factor)
	}
	return nil
}

// Whether to animate calendar swipes
// Should default to true
type AnimateCalendarSwipe struct {
	Enabled bool `json:"value"`
}

func (entry *AnimateCalendarSwipe) Key() string {
	return KeyAnimateCalendarSwipe
}
func (entry *AnimateCalendarSwipe) Default() {
	entry.Enabled = true
}
func (entry *AnimateCalendarSwipe) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *AnimateCalendarSwipe) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// Whether to animate small calendar swipes
// Should default to false
type AnimateSmallCalendarSwipe struct {
	Enabled bool `json:"value"`
}

func (entry *AnimateSmallCalendarSwipe) Key() string {
	return KeyAnimateSmallCalendarSwipe
}
func (entry *AnimateSmallCalendarSwipe) Default() {
	entry.Enabled = false
}
func (entry *AnimateSmallCalendarSwipe) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *AnimateSmallCalendarSwipe) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

// Whether to animate month selection swipes
// Should default to true
type AnimateMonthSelectionSwipe struct {
	Enabled bool `json:"value"`
}

func (entry *AnimateMonthSelectionSwipe) Key() string {
	return KeyAnimateMonthSelectionSwipe
}
func (entry *AnimateMonthSelectionSwipe) Default() {
	entry.Enabled = true
}
func (entry *AnimateMonthSelectionSwipe) MarshalJSON() ([]byte, error) {
	return common.MarshalBool(entry.Enabled), nil
}
func (entry *AnimateMonthSelectionSwipe) UnmarshalJSON(data []byte) (err error) {
	entry.Enabled, err = common.UnmarshalBool(data)
	return err
}

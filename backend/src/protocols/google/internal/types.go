package google

// https://developers.google.com/workspace/calendar/api/v3/reference/

type Calendar struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"summary"`
	Description string `json:"description,omitempty"`
}

type CalendarListEntry struct {
	Id              string `json:"id"`
	Name            string `json:"summary"`
	Description     string `json:"description,omitempty"`
	ColorId         string `json:"colorId,omitempty"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
	ForegroundColor string `json:"foregroundColor,omitempty"`
	Primary         bool   `json:"primary,omitempty"`
}

type ColorDefinition struct {
	Background string `json:"background"`
	Foreground string `json:"foreground"`
}

type Colors struct {
	Calendar map[string]ColorDefinition `json:"calendar"`
	Event    map[string]ColorDefinition `json:"event"`
}

type TimeDefinition struct {
	Date     string `json:"date,omitempty"`
	DateTime string `json:"dateTime,omitempty"`
	TimeZone string `json:"timeZone,omitempty"`
}

type Event struct {
	Id               string         `json:"id,omitempty"`
	Name             string         `json:"summary"`
	Description      string         `json:"description,omitempty"`
	ColorId          string         `json:"colorId,omitempty"`
	Start            TimeDefinition `json:"start"`
	End              TimeDefinition `json:"end"`
	Recurrence       []string       `json:"recurrence,omitempty"`
	IcalUid          string         `json:"icalUid,omitempty"`
	RecurringEventId string         `json:"recurringEventId,omitempty"`
}

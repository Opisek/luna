package google

// https://developers.google.com/workspace/calendar/api/v3/reference/

type CalendarListEntry struct {
	Id              string `json:"id"`
	Name            string `json:"summary"`
	Description     string `json:"description"`
	ColorId         string `json:"colorId"`
	BackgroundColor string `json:"backgroundColor"`
	ForegroundColor string `json:"foregroundColor"`
}

type ColorDefinition struct {
	Background string `json:"background"`
	Foreground string `json:"foreground"`
}

type Colors struct {
	Calendar map[string]ColorDefinition `json:"calendar"`
	Event    map[string]ColorDefinition `json:"event"`
}

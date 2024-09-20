package parsing

import (
	"encoding/json"
	"fmt"
	"luna-backend/db/internal/tables"
	"luna-backend/interface/primitives"
	"luna-backend/interface/protocols/caldav"
	"luna-backend/types"
)

func ParseEventSettings(sourceType string, settings []byte) (primitives.CalendarSettings, error) {
	switch sourceType {
	case types.SourceCaldav:
		parsedSettings := &caldav.CaldavEventSettings{}
		err := json.Unmarshal(settings, parsedSettings)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal caldav settings: %v", err)
		}
		return parsedSettings, nil
	case types.SourceIcal:
		fallthrough
	default:
		return nil, fmt.Errorf("unknown source type %v", sourceType)
	}
}

func ParseEventEntry(calendar primitives.Calendar, id types.ID, color []byte, settings []byte) (*tables.EventEntry, error) {
	parsedSettings, err := ParseEventSettings(calendar.GetSource().GetType(), settings)
	if err != nil {
		return nil, fmt.Errorf("could not parse calendar settings: %v", err)
	}

	return &tables.EventEntry{
		Id:       id,
		Calendar: calendar,
		Color:    types.ColorFromBytes(color),
		Settings: parsedSettings,
	}, nil
}

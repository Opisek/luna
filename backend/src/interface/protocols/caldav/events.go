package caldav

import (
	"luna-backend/interface/primitives"
	"luna-backend/types"
)

type CaldavEvent struct {
	id       types.ID
	name     string
	desc     string
	settings *CaldavEventSettings
	calendar *CaldavCalendar
}

type CaldavEventSettings struct {
}

func (event *CaldavEventSettings) GetBytes() []byte {
	return []byte{}
}

func (event *CaldavEvent) GetId() types.ID {
	return event.id
}

func (event *CaldavEvent) GetName() string {
	return event.name
}

func (event *CaldavEvent) GetDesc() string {
	return event.desc
}

func (event *CaldavEvent) GetCalendar() types.ID {
	return event.calendar.GetId()
}

func (event *CaldavEvent) GetSettings() primitives.CalendarSettings {
	return event.settings
}

func (event *CaldavEvent) GetColor() *types.Color {
	return types.ColorFromVals(50, 50, 50)
}

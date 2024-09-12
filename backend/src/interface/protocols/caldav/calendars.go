package caldav

import (
	"encoding/json"
	"luna-backend/auth"
	"luna-backend/crypto"
	"luna-backend/interface/primitives/calendars"
	"luna-backend/types"
	"time"

	"github.com/emersion/go-webdav/caldav"
)

type CaldavCalendar struct {
	name     string
	desc     string
	source   types.ID
	settings *CaldavCalendarSettings
	auth     auth.AuthMethod
	client   *caldav.Client
}

type CaldavCalendarSettings struct {
	Url *types.Url `json:"url"`
}

func (settings *CaldavCalendarSettings) GetBytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (calendar *CaldavCalendar) GetId() types.ID {
	return crypto.DeriveID(calendar.source, calendar.settings.Url.Path)
}

func (calendar *CaldavCalendar) GetName() string {
	return calendar.name
}

func (calendar *CaldavCalendar) GetDesc() string {
	return calendar.desc
}

func (calendar *CaldavCalendar) GetSource() types.ID {
	return calendar.source
}

func (calendar *CaldavCalendar) GetAuth() auth.AuthMethod {
	return calendar.auth
}

func (calendar *CaldavCalendar) GetSettings() calendars.CalendarSettings {
	return calendar.settings
}

func (calendar *CaldavCalendar) GetColor() *types.Color {
	return types.ColorFromVals(50, 50, 50)
}

func (calendar *CaldavCalendar) GetEvents(start time.Time, end time.Time) ([]*types.Event, error) {
	return nil, nil
}

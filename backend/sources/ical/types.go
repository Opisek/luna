package ical

import (
	"luna-backend/types"
	"net/url"

	ics "github.com/arran4/golang-ical"
)

type IcalSource struct {
	settings *IcalSettings
}

type IcalSettings struct {
	Url *url.URL
}

func calendarFromIcal(rawCalendar *ics.Calendar) *types.Calendar {
	properties := make(map[string]string)
	for _, prop := range rawCalendar.CalendarProperties {
		properties[prop.IANAToken] = prop.Value
	}

	return &types.Calendar{
		Name:  properties[string(ics.PropertyXWRCalName)],
		Desc:  properties[string(ics.PropertyXWRCalDesc)],
		Color: nil,
	}
}

package ical

import (
	"luna-backend/auth"
	"luna-backend/sources"
	"luna-backend/types"
	"net/url"

	ics "github.com/arran4/golang-ical"
)

type IcalSource struct {
	id       sources.SourceId
	settings *IcalSettings
}

type IcalSettings struct {
	Url  *url.URL
	Auth auth.AuthMethod
}

func (source *IcalSource) GetId() sources.SourceId {
	return source.id
}

func (source *IcalSource) calendarFromIcal(rawCalendar *ics.Calendar) *types.Calendar {
	properties := make(map[string]string)
	for _, prop := range rawCalendar.CalendarProperties {
		properties[prop.IANAToken] = prop.Value
	}

	return &types.Calendar{
		Source: source.GetId().String(),
		Id:     "", // for now we assume that one ical link corresponds to exactly one calendar
		Name:   properties[string(ics.PropertyXWRCalName)],
		Desc:   properties[string(ics.PropertyXWRCalDesc)],
		Color:  nil,
	}
}

func (source *IcalSource) eventFromIcal(rawEvent *ics.VEvent) *types.Event {
	properties := make(map[string]string)
	for _, prop := range rawEvent.Properties {
		properties[prop.IANAToken] = prop.Value
	}

	//startString := properties[string(ics.ComponentPropertyDtStart)]
	//endString := properties[string(ics.ComponentPropertyDtEnd)]
	//durationString := properties[string(ics.PropertyDuration)]

	//startTime, err := time.Parse(time.RFC3339, startString)
	//if err != nil {
	//	return nil
	//}
	//endTime, err := time.Parse(time.RFC3339, endString)

	//return &types.Event{
	//	Name: properties[string(ics.PropertySummary)],
	//	Start:
	//}

	return nil
}

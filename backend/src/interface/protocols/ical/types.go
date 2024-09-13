package ical

import (
	"luna-backend/auth"
	"luna-backend/types"
)

type IcalSource struct {
	id       types.ID
	settings *IcalSettings
}

type IcalSettings struct {
	Url  *types.Url
	Auth auth.AuthMethod
}

func (source *IcalSource) GetId() types.ID {
	return source.id
}

//func (source *IcalSource) calendarFromIcal(rawCalendar *ics.Calendar) *primitives.Calendar {
//	properties := make(map[string]string)
//	for _, prop := range rawCalendar.CalendarProperties {
//		properties[prop.IANAToken] = prop.Value
//	}
//
//	return &types.Calendar{
//		Source: source.GetId(),
//		Id:     types.EmptyId(), // TODO: placeholder
//		Path:   "",              // TODO: placeholder
//		Name:   properties[string(ics.PropertyXWRCalName)],
//		Desc:   properties[string(ics.PropertyXWRCalDesc)],
//		Color:  nil,
//	}
//}

//func (source *IcalSource) eventFromIcal(rawEvent *ics.VEvent) *types.Event {
//	properties := make(map[string]string)
//	for _, prop := range rawEvent.Properties {
//		properties[prop.IANAToken] = prop.Value
//	}
//
//	//startString := properties[string(ics.ComponentPropertyDtStart)]
//	//endString := properties[string(ics.ComponentPropertyDtEnd)]
//	//durationString := properties[string(ics.PropertyDuration)]
//
//	//startTime, err := time.Parse(time.RFC3339, startString)
//	//if err != nil {
//	//	return nil
//	//}
//	//endTime, err := time.Parse(time.RFC3339, endString)
//
//	//return &types.Event{
//	//	Name: properties[string(ics.PropertySummary)],
//	//	Start:
//	//}
//
//	return nil
//}

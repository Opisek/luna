package caldav

import (
	"fmt"
	"luna-backend/crypto"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"time"

	"github.com/emersion/go-webdav/caldav"
)

type CaldavEvent struct {
	uid      string
	name     string
	desc     string
	settings *CaldavEventSettings
	calendar *CaldavCalendar
	start    time.Time
	end      time.Time
}

type CaldavEventSettings struct {
}

func emptyEvent() *CaldavEvent {
	return &CaldavEvent{}
}

// TODO: proper parsing of start, end, duration, etc.
func eventFromCaldav(calendar *CaldavCalendar, obj *caldav.CalendarObject) (*CaldavEvent, error) {
	eventIndex := -1
	for i, child := range obj.Data.Children {
		if child.Name == "VEVENT" {
			eventIndex = i
			break
		}
	}
	if eventIndex == -1 {
		return nil, fmt.Errorf("could not find VEVENT in calendar object %v", obj.Path)
	}

	uid := obj.Data.Children[eventIndex].Props.Get("UID")
	summary := obj.Data.Children[eventIndex].Props.Get("SUMMARY")
	//dtstart := obj.Data.Children[0].Props.Get("DTSTART") // TODO: why does this not work?
	//tzid := dtstart.Params.Get("TZID")

	//fmt.Println("load loc")
	//location, err := time.LoadLocation(tzid)
	//if err != nil {
	//	panic(err) // TODO: error handling
	//}

	//fmt.Println("parse start time")
	//startTime, err := time.ParseInLocation("20060102T150405", dtstart.Value, location)
	//if err != nil {
	//	panic(err) // TODO: error handling
	//}

	return &CaldavEvent{
		uid:      uid.Value,
		name:     summary.Value,
		desc:     summary.Value,
		settings: &CaldavEventSettings{},
		calendar: calendar,
		start:    time.Now(),
		end:      time.Now().Add(time.Hour),
	}, nil
}

func (event *CaldavEventSettings) GetBytes() []byte {
	return []byte{}
}

func (event *CaldavEvent) GetId() types.ID {
	return crypto.DeriveID(event.calendar.GetId(), event.uid)
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

func (event *CaldavEvent) GetSettings() primitives.EventSettings {
	return event.settings
}

func (event *CaldavEvent) GetColor() *types.Color {
	return types.ColorFromVals(50, 50, 50)
}

func (event *CaldavEvent) GetStart() time.Time {
	return event.start
}

func (event *CaldavEvent) GetEnd() time.Time {
	return event.end
}

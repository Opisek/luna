package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
)

var DateFormatRfc5545 = "20060102"
var DateTimeFormatRfc5545 = "20060102T150405"
var DateTimeUtcFormatRfc5545 = "20060102T150405Z"
var FormatsRfc5545 = []string{DateFormatRfc5545, DateTimeFormatRfc5545, DateTimeUtcFormatRfc5545}

type EventDate struct {
	start           *time.Time
	end             *time.Time
	duration        *time.Duration
	specifyDuration bool
	allDay          bool
	timezone        *time.Location
	recurrence      *EventRecurrence
}

func NewEventDateFromEndTime(start *time.Time, end *time.Time, allDay bool, recurrence *EventRecurrence) *EventDate {
	_, offset := start.Zone()
	timezone := start.Location()
	if allDay {
		newStart := start.Add(time.Duration(offset) * time.Second).UTC()

		var newEnd time.Time
		if end.After(*start) {
			_, offset = end.Zone()
			newEnd = end.Add(time.Duration(offset) * time.Second).UTC()
		} else {
			newEnd = start.Add(24 * time.Hour)
		}

		start = &newStart
		end = &newEnd
	}

	return &EventDate{
		start:           start,
		end:             end,
		allDay:          allDay,
		recurrence:      recurrence,
		timezone:        timezone,
		specifyDuration: false,
	}
}

func NewEventDateFromDuration(start *time.Time, duration *time.Duration, allDay bool, recurrence *EventRecurrence) *EventDate {
	_, offset := start.Zone()
	timezone := start.Location()
	if allDay {
		newStart := start.Add(time.Duration(offset) * time.Second).UTC()
		start = &newStart
	}

	return &EventDate{
		start:           start,
		duration:        duration,
		allDay:          allDay,
		timezone:        timezone,
		recurrence:      recurrence,
		specifyDuration: true,
	}
}

func NewEventDateFromSingleDay(start *time.Time, recurrence *EventRecurrence) *EventDate {
	_, offset := start.Zone()
	timezone := start.Location()
	newStart := start.Add(time.Duration(offset) * time.Second).UTC()
	newEnd := newStart

	start = &newStart
	end := &newEnd

	return &EventDate{
		start:           start,
		end:             end,
		allDay:          true,
		timezone:        timezone,
		recurrence:      recurrence,
		specifyDuration: false,
	}
}

func (ed *EventDate) Start() *time.Time {
	return ed.start
}

func (ed *EventDate) SetStart(start *time.Time) {
	ed.start = start
}

func (ed *EventDate) End() *time.Time {
	if ed.specifyDuration {
		endTime := ed.start.Add(*ed.duration)
		return &endTime
	} else {
		return ed.end
	}
}

func (ed *EventDate) SetEnd(end *time.Time) {
	ed.end = end
	ed.specifyDuration = false
}

func (ed *EventDate) Duration() *time.Duration {
	if ed.specifyDuration {
		return ed.duration
	} else {
		duration := ed.end.Sub(*ed.start)
		return &duration
	}
}

func (ed *EventDate) SpecifyDuration() bool {
	return ed.specifyDuration
}

func (ed *EventDate) Timezone() *time.Location {
	return ed.timezone
}

func (ed *EventDate) SetTimezone(timezone *time.Location) {
	ed.timezone = timezone
}

func (ed *EventDate) Recurrence() *EventRecurrence {
	if ed.recurrence == nil {
		return EmptyEventRecurrence()
	}
	return ed.recurrence
}

func (ed *EventDate) SetRecurrence(recurrence *EventRecurrence) {
	ed.recurrence = recurrence
}

func (ed *EventDate) AllDay() bool {
	return ed.allDay
}

type eventDateMarshalEnd struct {
	Start      *time.Time       `json:"start"`
	End        *time.Time       `json:"end"`
	Recurrence *EventRecurrence `json:"recurrence,omitempty"`
	AllDay     bool             `json:"allDay"`
}

func (ed *EventDate) MarshalJSON() ([]byte, error) {
	var recurrence *EventRecurrence
	if ed.Recurrence().Repeats() {
		recurrence = ed.recurrence
	}
	if ed.specifyDuration {
		endTime := ed.start.Add(*ed.duration)
		return json.Marshal(eventDateMarshalEnd{
			Start:      ed.start,
			End:        &endTime,
			Recurrence: recurrence,
			AllDay:     ed.allDay,
		})
	} else {
		return json.Marshal(eventDateMarshalEnd{
			Start:      ed.start,
			End:        ed.end,
			Recurrence: recurrence,
			AllDay:     ed.allDay,
		})
	}
}

func (ed *EventDate) UnmarshalJSON(data []byte) error {
	var edme eventDateMarshalEnd
	err := json.Unmarshal(data, &edme)
	if err != nil {
		return err
	}
	ed.start = edme.Start
	ed.end = edme.End
	ed.recurrence = edme.Recurrence
	ed.specifyDuration = false
	return nil
}

func (ed *EventDate) Clone() *EventDate {
	if ed.specifyDuration {
		return NewEventDateFromDuration(ed.start, ed.duration, ed.allDay, ed.recurrence.Clone())
	} else {
		return NewEventDateFromEndTime(ed.start, ed.end, ed.allDay, ed.recurrence.Clone())
	}
}

// RFC-5545 3.3.10, 3.8.5.3
type EventRecurrence struct {
	ruleSet           *rrule.Set // TODO: change to rruleset?
	modifiedInstances []time.Time
	allDay            bool
	timezone          *time.Location
}

func (er *EventRecurrence) Clone() *EventRecurrence {
	if !er.Repeats() {
		return EmptyEventRecurrence()
	}

	ruleSet := *er.ruleSet

	return &EventRecurrence{
		ruleSet:           &ruleSet,
		modifiedInstances: er.modifiedInstances,
		allDay:            er.allDay,
		timezone:          er.timezone,
	}
}

func (er *EventRecurrence) Repeats() bool {
	return er.ruleSet != nil && (er.ruleSet.GetRRule() != nil || len(er.ruleSet.GetRDate()) != 0)
}

func (er *EventRecurrence) RuleSet() *rrule.Set {
	ruleSet := rrule.Set{}
	if er.ruleSet == nil {
		return &ruleSet
	}
	if er.ruleSet.GetRRule() != nil {
		ruleSet.RRule(er.ruleSet.GetRRule())
	}
	for _, date := range er.ruleSet.GetRDate() {
		ruleSet.RDate(date)
	}
	for _, date := range er.ruleSet.GetExDate() {
		ruleSet.ExDate(date)
	}
	return &ruleSet
}

func (er *EventRecurrence) EffectiveRuleSet() *rrule.Set {
	ruleSet := er.RuleSet()
	for _, modifiedInstance := range er.modifiedInstances {
		ruleSet.ExDate(modifiedInstance)
	}
	return ruleSet
}

func (er *EventRecurrence) AddException(exceptionTime *time.Time) {
	er.ruleSet.ExDate(*exceptionTime)
}

func (er *EventRecurrence) MarkModification(modifiedTime *time.Time) {
	er.modifiedInstances = append(er.modifiedInstances, *modifiedTime)
}

func EmptyEventRecurrence() *EventRecurrence {
	return &EventRecurrence{
		ruleSet: nil,
	}
}

func EventRecurrenceFromIcal(props *ical.Props) (*EventRecurrence, error) {
	var rset = rrule.Set{}

	roption, err := props.RecurrenceRule()
	if err != nil {
		return nil, fmt.Errorf("could not get recurrence rule: %v", err)
	}
	if roption != nil {
		rrule, err := rrule.NewRRule(*roption)
		if err != nil {
			return nil, fmt.Errorf("could not build recurrence rule: %v", err)
		}
		rset.RRule(rrule)
	}

	for _, prop := range props.Values(ical.PropExceptionDates) {
		exdateTime, _, err := ParseIcalTime(&prop)
		if err == nil {
			rset.ExDate(*exdateTime)
		}
	}

	for _, prop := range props.Values(ical.PropRecurrenceDates) {
		rdateTime, _, err := ParseIcalTime(&prop)
		if err == nil {
			rset.RDate(*rdateTime)
		}
	}

	var eventRecurrence *EventRecurrence
	if rset.GetRRule() == nil && len(rset.GetRDate()) == 0 {
		eventRecurrence = EmptyEventRecurrence()
	} else {
		var timezone *time.Location

		if dtstart := rset.GetDTStart(); dtstart.Unix() != 0 {
			timezone = dtstart.Location()
		} else if rdate := rset.GetRDate(); len(rdate) != 0 {
			timezone = rdate[0].Location()
		} else if exdate := rset.GetExDate(); len(exdate) != 0 {
			timezone = exdate[0].Location()
		}

		if timezone == nil {
			timezone = time.UTC
		}

		eventRecurrence = &EventRecurrence{
			ruleSet:  &rset,
			allDay:   true, // TODO
			timezone: timezone,
		}
	}

	return eventRecurrence, nil
}

func EventRecurrenceToIcal(recurrence *EventRecurrence) *ical.Props {
	if !recurrence.Repeats() {
		return nil
	}

	props := make(ical.Props)

	if rule := recurrence.RuleSet().GetRRule(); rule != nil {
		prop := ical.NewProp(ical.PropRecurrenceRule)
		prop.SetValueType(ical.ValueRecurrence)
		prop.Value = rule.Options.RRuleString()
		props.Set(prop)
	}
	if rdate := recurrence.RuleSet().GetRDate(); len(rdate) != 0 {
		for _, date := range rdate {
			prop := ical.NewProp(ical.PropRecurrenceDates)
			prop.SetValueType(ical.ValueDateTime)
			prop.Value = SerializeIcalTime(&date, recurrence.allDay, false)
			prop.Params.Add(ical.ParamTimezoneID, recurrence.timezone.String())
			if recurrence.allDay {
				prop.Params.Add(ical.ParamValue, "DATE")
			} else {
				prop.Params.Add(ical.ParamValue, "DATE-TIME")
			}
			props.Add(prop)
		}
	}
	if exdate := recurrence.RuleSet().GetExDate(); len(exdate) != 0 {
		for _, date := range exdate {
			prop := ical.NewProp(ical.PropExceptionDates)
			prop.SetValueType(ical.ValueDateTime)
			prop.Value = SerializeIcalTime(&date, recurrence.allDay, false)
			prop.Params.Add(ical.ParamTimezoneID, recurrence.timezone.String())
			if recurrence.allDay {
				prop.Params.Add(ical.ParamValue, "DATE")
			} else {
				prop.Params.Add(ical.ParamValue, "DATE-TIME")
			}
			props.Add(prop)
		}
	}

	return &props
}

func EventRecurrenceFromStrings(rrule string, rdate string, exdate string) (*EventRecurrence, error) {
	if rrule == "" && rdate == "" {
		return EmptyEventRecurrence(), nil
	}

	slice := make([]string, 0, 3)
	if rrule != "" {
		slice = append(slice, rrule)
	}
	if rdate != "" {
		slice = append(slice, rdate)
	}
	if exdate != "" {
		slice = append(slice, exdate)
	}

	return EventRecurrenceFromLines(slice)
}

func EventRecurrenceFromLines(lines []string) (*EventRecurrence, error) {
	if len(lines) == 0 {
		return EmptyEventRecurrence(), nil
	}

	rset, err := rrule.StrToRRuleSet(strings.Join(lines, "\n"))
	if err != nil {
		return nil, fmt.Errorf("could not get recurrence rule: %v", err)
	}

	if rset == nil {
		return EmptyEventRecurrence(), nil
	}

	return &EventRecurrence{
		ruleSet:  rset,
		allDay:   true, // TODO
		timezone: nil,  // TODO
	}, nil
}

func ParseIcalTime(icalTime *ical.Prop) (*time.Time, *time.Location, error) {
	if icalTime == nil || icalTime.Value == "" {
		return nil, nil, fmt.Errorf("time property is nil or empty")
	}

	var tzid string
	if icalTime.Value[len(icalTime.Value)-1] == 'Z' {
		tzid = "UTC"
	} else {
		tzidParam := icalTime.Params.Get("TZID")
		if tzidParam == "" {
			tzid = "Local"
		} else {
			tzid = tzidParam
		}
	}

	location, err := time.LoadLocation(tzid)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse timezone location %v: %v", tzid, err)
	}

	for _, format := range FormatsRfc5545 {
		parsedTime, err := time.ParseInLocation(format, icalTime.Value, location)
		if err == nil {
			return &parsedTime, location, nil
		}
	}

	return nil, nil, fmt.Errorf("could not parse timestamp %v", icalTime.Value)
}

func SerializeIcalTime(time *time.Time, allDay bool, utc bool) string {
	if allDay {
		return time.Format(DateFormatRfc5545)
	} else if utc {
		return time.UTC().Format(DateTimeUtcFormatRfc5545)
	} else {
		return time.Format(DateTimeFormatRfc5545)
	}
}

func SerializeIcalTimeSlice(slice []time.Time, allDay bool) []string {
	res := make([]string, len(slice))
	for i, time := range slice {
		res[i] = SerializeIcalTime(&time, allDay, false)
	}
	return res
}

type eventRecurrenceMarshal struct {
	Rrule  string `json:"RRULE,omitempty"`
	Exdate string `json:"EXDATE,omitempty"`
	Rdate  string `json:"RDATE,omitempty"`
}

func (er EventRecurrence) RruleString() string {
	if !er.Repeats() || er.ruleSet.GetRRule() == nil || er.ruleSet.GetRRule().String() == "" {
		return ""
	}
	return fmt.Sprintf("%s", er.ruleSet.GetRRule().Options.RRuleString())
}

func (er EventRecurrence) RdateString() string {
	if !er.Repeats() {
		return ""
	}
	var valueType string
	if er.allDay {
		valueType = "DATE"
	} else {
		valueType = "DATE-TIME"
	}
	return fmt.Sprintf("RDATE;VALUE=%s;TZID=%s:%s", valueType, er.timezone.String(), strings.Join(SerializeIcalTimeSlice(er.ruleSet.GetRDate(), er.allDay), ","))
}

func (er EventRecurrence) ExdateString() string {
	if !er.Repeats() {
		return ""
	}
	var valueType string
	if er.allDay {
		valueType = "DATE"
	} else {
		valueType = "DATE-TIME"
	}
	return fmt.Sprintf("EXDATE;VALUE=%s;TZID=%s:%s", valueType, er.timezone.String(), strings.Join(SerializeIcalTimeSlice(er.ruleSet.GetExDate(), er.allDay), ","))
}

func (er EventRecurrence) Lines() []string {
	if !er.Repeats() {
		return []string{}
	}

	var lines = make([]string, 0, 3)
	if str := er.RruleString(); str != "" {
		lines = append(lines, str)
	}
	if str := er.RdateString(); str != "" {
		lines = append(lines, str)
	}
	if str := er.ExdateString(); str != "" {
		lines = append(lines, str)
	}

	return lines
}

func (er EventRecurrence) MarshalJSON() ([]byte, error) {
	if er.Repeats() {
		casted := eventRecurrenceMarshal{
			Rrule:  er.RruleString(),
			Rdate:  er.RdateString(),
			Exdate: er.ExdateString(),
		}
		return json.Marshal(casted)
	} else {
		return []byte("null"), nil
	}
}

func (er *EventRecurrence) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	if str == "false" || str == "null" || str == "" {
		er = EmptyEventRecurrence()
		return nil
	}

	var unmarshaled eventRecurrenceMarshal
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return
	}

	er, err = EventRecurrenceFromStrings(unmarshaled.Rrule, unmarshaled.Rdate, unmarshaled.Exdate)
	return
}

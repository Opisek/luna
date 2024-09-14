package types

import (
	"encoding/json"
	"time"
)

type EventDate struct {
	start           time.Time
	end             time.Time
	duration        time.Duration
	specifyDuration bool
	allDay          bool
	recurrence      EventRecurrence
}

type eventDateMarshalEnd struct {
	Start      time.Time       `json:"start"`
	End        time.Time       `json:"end"`
	AllDay     bool            `json:"allDay"`
	Recurrence EventRecurrence `json:"recurrence"`
}

type eventDateMarshalDuration struct {
	Start      time.Time       `json:"start"`
	Duration   time.Duration   `json:"duration"`
	AllDay     bool            `json:"allDay"`
	Recurrence EventRecurrence `json:"recurrence"`
}

func (ed *EventDate) MarshalJSON() ([]byte, error) {
	if ed.specifyDuration {
		return json.Marshal(eventDateMarshalDuration{
			Start:      ed.start,
			Duration:   ed.duration,
			AllDay:     ed.allDay,
			Recurrence: ed.recurrence,
		})
	} else {
		return json.Marshal(eventDateMarshalEnd{
			Start:      ed.start,
			End:        ed.end,
			AllDay:     ed.allDay,
			Recurrence: ed.recurrence,
		})
	}
}

func (ed *EventDate) UnmarshalJSON(data []byte) error {
	var edme eventDateMarshalEnd
	if err := json.Unmarshal(data, &edme); err == nil {
		ed.start = edme.Start
		ed.end = edme.End
		ed.allDay = edme.AllDay
		ed.recurrence = edme.Recurrence
		ed.specifyDuration = false
		return nil
	}

	var edmd eventDateMarshalDuration
	if err := json.Unmarshal(data, &edmd); err == nil {
		ed.start = edmd.Start
		ed.duration = edmd.Duration
		ed.allDay = edmd.AllDay
		ed.recurrence = edmd.Recurrence
		ed.specifyDuration = true
		return nil
	}

	return nil
}

// TODO: validation according to RFC-5545 3.3.10, 3.8.5.3
type EventRecurrence map[string]string

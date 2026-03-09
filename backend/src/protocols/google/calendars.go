package google

import (
	"encoding/json"
	"luna-backend/crypto"
	"luna-backend/errors"
	google "luna-backend/protocols/google/internal"
	"luna-backend/types"
	"net/http"
	"time"
)

type GoogleCalendar struct {
	name       string
	desc       string
	color      *types.Color
	overridden bool
	settings   *GoogleCalendarSettings
	source     *GoogleSource
}

type GoogleCalendarSettings struct {
	InternalId string `json:"google_id"`
}

func (source *GoogleSource) calendarFromGoogle(calListEntry *google.CalendarListEntry, q types.DatabaseQueries) (*GoogleCalendar, *errors.ErrorTrace) {
	var tr *errors.ErrorTrace

	var col *types.Color
	if calListEntry.BackgroundColor == "" {
		col, _, tr = source.getColorById(calListEntry.ColorId, true, q)
		if tr != nil {
			return nil, tr.
				Append(errors.LvlDebug, "Could not parse calendar %v", calListEntry.Id).
				AltStr(errors.LvlWordy, "Could not parse calendar")
		}
	} else {
		var err error
		col, err = types.ParseColor(calListEntry.BackgroundColor)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				Append(errors.LvlDebug, "Could not parse color %v", calListEntry.BackgroundColor).
				Append(errors.LvlDebug, "Could not parse calendar %v", calListEntry.Id).
				AltStr(errors.LvlWordy, "Could not parse calendar")
		}
	}

	settings := &GoogleCalendarSettings{
		InternalId: calListEntry.Id,
	}

	calendar := &GoogleCalendar{
		name:       calListEntry.Name,
		desc:       calListEntry.Description,
		color:      col,
		overridden: false,
		settings:   settings,
		source:     source,
	}

	return calendar, nil
}

func (settings *GoogleCalendarSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genCalId(sourceId types.ID, path string) types.ID {
	return crypto.DeriveID(sourceId, path)
}

func (calendar *GoogleCalendar) GetId() types.ID {
	return genCalId(calendar.source.id, calendar.settings.InternalId)
}

func (calendar *GoogleCalendar) GetName() string {
	return calendar.name
}

func (calendar *GoogleCalendar) SetName(name string) {
	calendar.name = name
}

func (calendar *GoogleCalendar) GetDesc() string {
	return calendar.desc
}

func (calendar *GoogleCalendar) SetDesc(desc string) {
	calendar.desc = desc
}

func (calendar *GoogleCalendar) GetSource() types.Source {
	return calendar.source
}

func (calendar *GoogleCalendar) GetSettings() types.CalendarSettings {
	return calendar.settings
}

func (calendar *GoogleCalendar) GetColor() *types.Color {
	if calendar.color == nil {
		return types.ColorEmpty
	} else {
		return calendar.color
	}
}

func (calendar *GoogleCalendar) SetColor(color *types.Color) {
	calendar.color = color
}

func (calendar *GoogleCalendar) GetOverridden() bool {
	return calendar.overridden
}

func (calendar *GoogleCalendar) SetOverridden(overridden bool) {
	calendar.overridden = overridden
}

func (calendar *GoogleCalendar) GetEvents(start time.Time, end time.Time, q types.DatabaseQueries) ([]types.Event, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusNotImplemented)
}

func (calendar *GoogleCalendar) GetEvent(settings types.EventSettings, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusNotImplemented)
}

func (calendar *GoogleCalendar) AddEvent(name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusNotImplemented)
}

func (calendar *GoogleCalendar) EditEvent(originalEvent types.Event, name string, desc string, color *types.Color, date *types.EventDate, _ bool, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	return nil, errors.New().Status(http.StatusNotImplemented)
}

func (calendar *GoogleCalendar) DeleteEvent(event types.Event, q types.DatabaseQueries) *errors.ErrorTrace {
	return errors.New().Status(http.StatusNotImplemented)
}

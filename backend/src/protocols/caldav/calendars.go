package caldav

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/errors"
	supplementary_caldav "luna-backend/protocols/caldav/internal"
	common "luna-backend/protocols/internal"
	"luna-backend/types"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

type CaldavCalendar struct {
	name       string
	desc       string
	color      *types.Color
	overridden bool
	settings   *CaldavCalendarSettings
	source     *CaldavSource
	client     *caldav.Client
}

type CaldavCalendarSettings struct {
	Url         *types.Url      `json:"url"`
	rawCalendar caldav.Calendar `json:"-"`
}

func (source *CaldavSource) calendarFromCaldav(rawCalendar caldav.Calendar) (*CaldavCalendar, *errors.ErrorTrace) {
	url, err := types.NewUrl(rawCalendar.Path)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse URL %v", rawCalendar.Path).
			Append(errors.LvlWordy, "Could not parse calendar")
	}

	props, tr := supplementary_caldav.PropFind(source.settings.Url, url, []string{"I:calendar-color"}, source.auth, source.ctx)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse calendar")
	}

	var color *types.Color = nil
	colProp, exists := props["I:calendar-color"]
	if exists && colProp.Found {
		color, err = types.ParseColor(colProp.Value)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not parse calendar color %v", colProp.Value).
				Append(errors.LvlWordy, "Could not parse calendar")
		}
	}

	settings := &CaldavCalendarSettings{
		Url:         url,
		rawCalendar: rawCalendar,
	}

	calendar := &CaldavCalendar{
		name:       rawCalendar.Name,
		desc:       rawCalendar.Description,
		color:      color,
		overridden: false,
		settings:   settings,
		source:     source,
		client:     source.client,
	}

	return calendar, nil
}

func (settings *CaldavCalendarSettings) Bytes() []byte {
	bytes, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	return bytes
}

func genCalId(sourceId types.ID, path string) types.ID {
	return crypto.DeriveID(sourceId, path)
}

func (calendar *CaldavCalendar) GetId() types.ID {
	return genCalId(calendar.source.id, calendar.settings.Url.Path)
}

func (calendar *CaldavCalendar) GetName() string {
	return calendar.name
}

func (calendar *CaldavCalendar) SetName(name string) {
	calendar.name = name
}

func (calendar *CaldavCalendar) GetDesc() string {
	return calendar.desc
}

func (calendar *CaldavCalendar) SetDesc(desc string) {
	calendar.desc = desc
}

func (calendar *CaldavCalendar) GetSource() types.Source {
	return calendar.source
}

func (calendar *CaldavCalendar) GetSettings() types.CalendarSettings {
	return calendar.settings
}

func (calendar *CaldavCalendar) GetColor() *types.Color {
	if calendar.color == nil {
		return types.ColorEmpty
	} else {
		return calendar.color
	}
}

func (calendar *CaldavCalendar) SetColor(color *types.Color) {
	calendar.color = color
}

func (calendar *CaldavCalendar) GetOverridden() bool {
	return calendar.overridden
}

func (calendar *CaldavCalendar) SetOverridden(overridden bool) {
	calendar.overridden = overridden
}

func (calendar *CaldavCalendar) CanEdit() bool {
	return true
}

func (calendar *CaldavCalendar) CanDelete() bool {
	return true
}

func (calendar *CaldavCalendar) CanAddEvents() bool {
	return true
}

func (calendar *CaldavCalendar) convertEvent(event *caldav.CalendarObject, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	convertedEvent, err := calendar.eventFromCaldav(event, q)
	if err != nil {
		return nil, err.
			Append(errors.LvlDebug, "Could not convert calendar %v", event.Path).
			AltStr(errors.LvlWordy, "Could not convert calendar")
	}

	castedEvent := (types.Event)(convertedEvent)

	return castedEvent, nil
}

func normalizeETag(etag string) (string, error) {
	etag = strings.TrimSpace(etag)
	if etag == "" {
		return "", nil
	}
	// strip the weak etag specifier and any whitespace
	if strings.HasPrefix(etag, "W/") {
		etag = strings.TrimSpace(etag[2:])
	}

	// make sure the etag is actually a valid format
	if len(etag) < 2 || etag[0] != '"' || etag[len(etag)-1] != '"' {
		return "", fmt.Errorf("caldav: invalid ETag %q", etag)
	}

	return etag[1 : len(etag)-1], nil
}

func (calendar *CaldavCalendar) getEvents(query *caldav.CalendarQuery, q types.DatabaseQueries) ([]types.Event, *errors.ErrorTrace) {
	client, tr := calendar.source.getClient()
	if tr != nil {
		return nil, tr.
			Append(errors.LvlBroad, "Could not get events")
	}

	events, err := client.QueryCalendar(q.GetContext(), calendar.settings.Url.String(), query)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not get events")
	}

	masterEvents := make(map[string]int)
	masterEventIndices := make(map[int]bool)

	convertedEvents := make([]types.Event, len(events))
	for i, event := range events {
		convertedEvents[i], tr = calendar.convertEvent(&event, q)
		if tr != nil {
			return nil, tr.
				Append(errors.LvlBroad, "Could not get events")
		}

		eventSettings := convertedEvents[i].GetSettings().(*CaldavEventSettings)

		if eventSettings.RecurrenceId == "" {
			masterEvents[eventSettings.Uid] = i
			masterEventIndices[i] = true
		}
	}

	// Internally note all the modified recurrence instances for each master event so that we don't expand these later
	for i, event := range convertedEvents {
		if masterEventIndices[i] {
			continue
		}
		eventSettings := event.GetSettings().(*CaldavEventSettings)
		if masterEvent, exists := masterEvents[eventSettings.Uid]; exists {
			convertedEvents[masterEvent].GetDate().Recurrence().AddModifiedInstance(common.ExtractDateFromRecurrenceId(event))
		}
	}

	return convertedEvents, nil
}

func (calendar *CaldavCalendar) getCalendarObjectFallback(path string, q types.DatabaseQueries) (*caldav.CalendarObject, *errors.ErrorTrace) {
	target := *calendar.source.settings.Url.URL()
	target.Path = path
	target.RawQuery = ""
	target.Fragment = ""

	req, err := http.NewRequestWithContext(q.GetContext(), http.MethodGet, target.String(), nil)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not construct calendar data request")
	}
	req.Header.Set("Accept", ical.MIMEType)

	resp, err := calendar.source.auth.HttpClient().Do(req)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Call to get calendar data failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.InterpretRemoteError(errors.New().
			AddErr(errors.LvlDebug, fmt.Errorf("GET %q returned %s",path,resp.Status)), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not get calendar data")
	}

	mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not get content-type from calendar data")
	}
	if !strings.EqualFold(mediaType, ical.MIMEType) {
		return nil, errors.InterpretRemoteError(errors.New().
			AddErr(errors.LvlDebug, fmt.Errorf("Expected Content-Type %q, got %q", ical.MIMEType, mediaType)), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not get calendar data")
	}

	data, err := ical.NewDecoder(resp.Body).Decode()
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not decode calendar data")
	}

	obj := &caldav.CalendarObject{
		Path: resp.Request.URL.Path,
		Data: data,
	}

	// process/format the etag
	obj.ETag, err = normalizeETag(resp.Header.Get("ETag"))
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not find ETag in response header")
	}

	// handle other header info per the go-webdav
	if location := resp.Header.Get("Location"); location != "" {
		locationURL, err := url.Parse(location)
		if err != nil {
			return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
				Append(errors.LvlBroad, "Could not find Location in response header")
		}
		obj.Path = locationURL.Path
	}

	if value := resp.Header.Get("Content-Length"); value != "" {
		obj.ContentLength, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
				Append(errors.LvlBroad, "Could not find Content-Length in response header")
		}
	}

	if value := resp.Header.Get("Last-Modified"); value != "" {
		obj.ModTime, err = http.ParseTime(value)
		if err != nil {
			return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
				Append(errors.LvlBroad, "Could not find Last-Modified in response header")
		}
	}

	return obj, nil
}

func (calendar *CaldavCalendar) getCalendarObjectCompat(q types.DatabaseQueries, path string) (*caldav.CalendarObject, *errors.ErrorTrace) {
	// first make a default call
	obj, err := calendar.client.GetCalendarObject(q.GetContext(), path)
	if err == nil {
		return obj, nil
	}
	if !stderrors.Is(err, strconv.ErrSyntax) {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not get calendar data")
	}
	// try with a custom construction to handle weak ETags (non compliant to RFC 4791, but valid HTTP)
	obj, fallbackErr := calendar.getCalendarObjectFallback(path, q)
	if fallbackErr != nil {
		return nil, fallbackErr.Append(errors.LvlDebug, fmt.Sprintf("Standard GetCalendarObject failed: %v", err))
	}

	return obj, nil
}

func (calendar *CaldavCalendar) GetEvents(start time.Time, end time.Time, q types.DatabaseQueries) ([]types.Event, *errors.ErrorTrace) {
	return calendar.getEvents(&caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name: "VCALENDAR",
			Comps: []caldav.CalendarCompRequest{{
				Name: "VEVENT",
				Props: []string{
					"SUMMARY",
					"UID",
					"DTSTART",
					"DTEND",
					"DURATION",
				},
			}},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{{
				Name:  "VEVENT",
				Start: start,
				End:   end,
			}},
		},
	}, q)
}

func (calendar *CaldavCalendar) GetEvent(settings types.EventSettings, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	caldavSettings := settings.(*CaldavEventSettings)

	obj, err := calendar.getCalendarObjectCompat(q, caldavSettings.Url.Path)
	if err != nil {
		return nil, err.Append(errors.LvlBroad, "Could not get event")
	}

	cal, tr := calendar.convertEvent(obj, q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlBroad, "Could not get event")
	}

	return cal, nil
}

func setEventProps(cal *ical.Calendar, id string, name string, desc string, color *types.Color, date *types.EventDate) *errors.ErrorTrace {
	var event *ical.Event = nil
	for _, child := range cal.Children {
		if child.Name == "VEVENT" {
			event = ical.NewEvent()
			event.Component = child
			break
		}
	}
	if event == nil {
		event = ical.NewEvent()
		cal.Children = append(cal.Children, event.Component)
	}

	event.Props.SetText(ical.PropUID, id)

	event.Props.SetText(ical.PropSummary, common.EscapeIcalString(name))

	if desc != "" {
		event.Props.SetText(ical.PropDescription, common.EscapeIcalString(desc))
	} else {
		event.Props.Del(ical.PropDescription)
	}

	if color.IsEmpty() {
		event.Props.Del(ical.PropColor)
		event.Props.Del(common.PropColor)
		event.Props.Del(common.PropLastColorName)
	} else {
		colorName, exact := types.ColorToName(color)

		// According to the specification, the "COLOR" property must be a named CSS color.
		// To ensure compatibility, we map colors to the closest named CSS color for other clients,
		// and use a custom property for the exac color displayed in Luna.

		event.Props.SetText(ical.PropColor, colorName)
		if exact {
			event.Props.Del(common.PropColor)
			event.Props.Del(common.PropLastColorName)
		} else {
			event.Props.SetText(common.PropColor, color.String())
			// To detect when the color is changed by another client, we store the last color name in a custom property.
			event.Props.SetText(common.PropLastColorName, colorName)
		}
	}

	if date.AllDay() {
		event.Props.SetDate(ical.PropDateTimeStart, *date.Start())
	} else {
		event.Props.SetDateTime(ical.PropDateTimeStart, *date.Start())
	}
	if date.SpecifyDuration() {
		// TODO: figure this out

		return errors.New().Status(http.StatusNotImplemented)
		//event.Props.SetText(ical.PropDuration, *date.Duration())
	} else {
		if date.AllDay() {
			event.Props.SetDate(ical.PropDateTimeEnd, *date.End())
		} else {
			event.Props.SetDateTime(ical.PropDateTimeEnd, *date.End())
		}
		event.Props.Del(ical.PropDuration)
	}

	timestamp := time.Now().UTC()
	event.Props.SetDateTime(ical.PropDateTimeStamp, timestamp)
	//event.Props.SetDateTime(util.PropTimestamp, timestamp)

	cal.Props.SetText(ical.PropProductID, "Luna 0.1.0") // TODO: take version from common config

	cal.Props.SetText(ical.PropVersion, "2.0") // iCalendar version

	return nil
}

func (calendar *CaldavCalendar) AddEvent(name string, desc string, color *types.Color, date *types.EventDate, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	id := types.RandomId()
	cal := ical.NewCalendar()

	tr := setEventProps(cal, id.String(), name, desc, color, date)
	if tr != nil {
		return nil, tr.Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Could not set iCal properties").
			AltStr(errors.LvlPlain, "Malformed settings").
			Append(errors.LvlBroad, "Could not add event")
	}

	path := fmt.Sprintf("%v%v.ics", calendar.settings.Url.Path, id.String())

	_, err := calendar.client.PutCalendarObject(q.GetContext(), path, cal)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlBroad, "Could not add event")
	}

	obj, getCalErr := calendar.getCalendarObjectCompat(q, path)
	if getCalErr != nil {
		return nil, getCalErr.
			Append(errors.LvlWordy, "Could not get finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	finishedEvent, tr := calendar.eventFromCaldav(obj, q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse finished event").
			Append(errors.LvlBroad, "Could not add event")
	}

	return finishedEvent, nil
}

func (calendar *CaldavCalendar) EditEvent(originalEvent types.Event, name string, desc string, color *types.Color, date *types.EventDate, _ bool, q types.DatabaseQueries) (types.Event, *errors.ErrorTrace) {
	originalCaldavEvent := originalEvent.(*CaldavEvent)
	uid := originalCaldavEvent.GetSettings().(*CaldavEventSettings).Uid
	originalRawEvent := originalCaldavEvent.settings.rawEvent
	cal := originalRawEvent.Data

	tr := setEventProps(cal, uid, name, desc, color, date)
	if tr != nil {
		return nil, tr.Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Could not set iCal properties").
			AltStr(errors.LvlPlain, "Malformed settings").
			Append(errors.LvlBroad, "Could not add event")
	}

	_, err := calendar.client.PutCalendarObject(q.GetContext(), originalRawEvent.Path, cal)
	if err != nil {
		return nil, errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "calendar", "CalDAV calendar").
			Append(errors.LvlWordy, "Could not edit event").
			AltStr(errors.LvlBroad, "Could not edit event")
	}

	obj, getCalErr := calendar.getCalendarObjectCompat(q, originalRawEvent.Path)
	if getCalErr != nil {
		return nil, getCalErr.
			Append(errors.LvlWordy, "Could not get finished event").
			Append(errors.LvlBroad, "Could not edit event")
	}

	finishedEvent, tr := calendar.eventFromCaldav(obj, q)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not parse finished event").
			Append(errors.LvlBroad, "Could not edit event")
	}

	return finishedEvent, nil
}

func (calendar *CaldavCalendar) DeleteEvent(event types.Event, q types.DatabaseQueries) *errors.ErrorTrace {
	settings := event.GetSettings().(*CaldavEventSettings)

	err := calendar.client.RemoveAll(q.GetContext(), settings.Url.Path)
	if err != nil {
		return errors.InterpretRemoteError(errors.New().AddErr(errors.LvlDebug, err), "event", "CalDAV event").
			Append(errors.LvlBroad, "Could not delete event")
	}

	return nil
}

func (calendar *CaldavCalendar) SupplyContext(ctx context.Context) {
	calendar.source.SupplyContext(ctx)
}

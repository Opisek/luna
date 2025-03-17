package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type exposedEvent struct {
	Id       types.ID         `json:"id"`
	Calendar types.ID         `json:"calendar"`
	Name     string           `json:"name"`
	Desc     string           `json:"desc"`
	Color    *types.Color     `json:"color"`
	Date     *types.EventDate `json:"date"`
}

func GetEvents(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	calendarId, tr := util.GetId(c, "calendar")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Get the requested calendar
	calendar, tr := u.Tx.Queries().GetCalendar(userId, calendarId)
	if tr != nil {
		u.Error(tr)
		return
	}

	// Get the associated events
	startStr := c.Query("start")
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		u.Error(errors.New().
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Missing or malformed start time"))
		return
	}
	endStr := c.Query("end")
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		u.Error(errors.New().
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Missing or malformed end time"))
		return
	}
	if startTime.After(endTime) {
		u.Error(errors.New().
			Append(errors.LvlPlain, "Start time must not be after end time"))
		return
	}
	if endTime.Sub(startTime) > time.Hour*24*365 {
		endTime = startTime.Add(time.Hour * 24 * 365)
	}

	eventsFromCal, tr := calendar.GetEvents(startTime, endTime, u.Tx.Queries())
	if tr != nil {
		u.Error(tr)
		return
	}

	events, tr := u.Tx.Queries().ReconcileEvents(eventsFromCal)
	if tr != nil {
		u.Error(tr)
		return
	}

	// Expand recurring events
	expandedEvents := make([]primitives.Event, len(events))
	count := 0
	for _, event := range events {
		expanded, tr := primitives.ExpandRecurrence(event, &startTime, &endTime)
		if tr != nil {
			u.Error(tr)
			return
		}

		if len(expanded) > 1 {
			newRes := make([]primitives.Event, len(expandedEvents)-1+len(expanded))
			copy(newRes, expandedEvents[:count])
			expandedEvents = newRes
		}

		for _, e := range expanded {
			expandedEvents[count] = e
			count++
		}
	}

	// Convert to exposed format
	convertedEvents := make([]exposedEvent, count)
	for i, event := range expandedEvents[:count] {
		if event.GetName() == "" { // TODO: error handling
			continue
		}

		convertedEvents[i] = exposedEvent{
			Id:       event.GetId(),
			Calendar: event.GetCalendar().GetId(),
			Name:     event.GetName(),
			Desc:     event.GetDesc(),
			Color:    event.GetColor(),
			Date:     event.GetDate(),
		}
	}

	u.Success(&gin.H{"events": convertedEvents})
}

func GetEvent(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	eventId, tr := util.GetId(c, "event")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Get event
	event, err := u.Tx.Queries().GetEvent(userId, eventId)
	if err != nil {
		u.Error(err)
		return
	}

	// Convert to exposed format
	convertedCal := exposedEvent{
		Id:       event.GetId(),
		Calendar: event.GetCalendar().GetId(),
		Name:     event.GetName(),
		Desc:     event.GetDesc(),
		Color:    event.GetColor(),
		Date:     event.GetDate(),
		//Settings: event.GetSettings(),
	}

	u.Success(&gin.H{"event": convertedCal})
}

func PutEvent(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	calendarId, tr := util.GetId(c, "calendar")
	if tr != nil {
		u.Error(tr)
		return
	}

	calendar, tr := u.Tx.Queries().GetCalendar(userId, calendarId)
	if tr != nil {
		u.Error(tr)
		return
	}

	eventName := c.PostForm("name")
	if eventName == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing name"))
		return
	}

	eventDesc := c.PostForm("desc")

	eventColor, err := types.ParseColor(c.PostForm("color"))
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Missing or malformed color"))
		return
	}

	eventDateAllDay := c.PostForm("date_all_day") == "true"

	eventDateStartStr := c.PostForm("date_start")
	eventDateStart, err := time.Parse(time.RFC3339, eventDateStartStr)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Missing or malformed start time"))
		return
	}

	eventDateEndStr := c.PostForm("date_end")
	eventDateDurationStr := c.PostForm("date_duration")

	eventDateEnd, endErr := time.Parse(time.RFC3339, eventDateEndStr)
	eventDateDuration, durationErr := time.ParseDuration(eventDateDurationStr)

	var date *types.EventDate
	if (endErr != nil && durationErr != nil) || (endErr == nil && durationErr == nil) {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, endErr).AndErr(durationErr).
			Append(errors.LvlPlain, "Missing or malformed date end or duration"))
		return
	} else if endErr == nil && durationErr == nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Cannot specify both end and duration"))
		return
	} else if endErr == nil {
		date = types.NewEventDateFromEndTime(&eventDateStart, &eventDateEnd, eventDateAllDay, nil)
	} else {
		date = types.NewEventDateFromDuration(&eventDateStart, &eventDateDuration, eventDateAllDay, nil)
	}

	event, tr := calendar.AddEvent(eventName, eventDesc, eventColor, date, u.Tx.Queries())
	if tr != nil {
		u.Error(tr)
		return
	}

	tr = u.Tx.Queries().InsertEvent(event)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{"id": event.GetId().String()})
}

func PatchEvent(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	eventId, err := util.GetId(c, "event")
	if err != nil {
		u.Error(err)
		return
	}

	event, err := u.Tx.Queries().GetEvent(userId, eventId)
	if err != nil {
		u.Error(err)
		return
	}

	newEventName := c.PostForm("name")

	newEventDesc := c.PostForm("desc")

	newEventColor, colErr := types.ParseColor(c.PostForm("color"))

	eventDateAllDay := c.PostForm("date_all_day") == "true"

	eventDateStartStr := c.PostForm("date_start")
	eventDateStart, startErr := time.Parse(time.RFC3339, eventDateStartStr)

	eventDateEndStr := c.PostForm("date_end")
	eventDateDurationStr := c.PostForm("date_duration")

	eventDateEnd, endErr := time.Parse(time.RFC3339, eventDateEndStr)
	eventDateDuration, durationErr := time.ParseDuration(eventDateDurationStr)

	if newEventName == "" && newEventDesc == event.GetDesc() && colErr != nil && startErr != nil && endErr != nil && durationErr != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Nothing to change"))
		return
	}

	if newEventName == "" {
		newEventName = event.GetName()
	}

	if colErr != nil {
		newEventColor = event.GetColor()
	}

	var newEventDate *types.EventDate
	if startErr != nil && endErr != nil && durationErr != nil {
		newEventDate = event.GetDate()
	} else {
		if startErr != nil {
			eventDateStart = *event.GetDate().Start()
		}
		if endErr != nil && durationErr == nil {
			if event.GetDate().SpecifyDuration() {
				eventDateDuration = *event.GetDate().Duration()
				newEventDate = types.NewEventDateFromDuration(&eventDateStart, &eventDateDuration, eventDateAllDay, nil)
			} else {
				eventDateEnd = *event.GetDate().End()
				newEventDate = types.NewEventDateFromEndTime(&eventDateStart, &eventDateEnd, eventDateAllDay, nil)
			}
		} else if endErr == nil && durationErr == nil {
			u.Error(errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Cannot specify both end and duration"))
			return
		} else if endErr == nil {
			newEventDate = types.NewEventDateFromEndTime(&eventDateStart, &eventDateEnd, eventDateAllDay, nil)
		} else {
			newEventDate = types.NewEventDateFromDuration(&eventDateStart, &eventDateDuration, eventDateAllDay, nil)
		}
	}

	newEvent, err := event.GetCalendar().EditEvent(event, newEventName, newEventDesc, newEventColor, newEventDate, u.Tx.Queries())
	if err != nil {
		u.Error(err)
		return
	}

	err = u.Tx.Queries().UpdateEvent(newEvent)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

func DeleteEvent(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	eventId, err := util.GetId(c, "event")
	if err != nil {
		u.Error(err)
		return
	}

	// Get event first
	event, err := u.Tx.Queries().GetEvent(userId, eventId)
	if err != nil {
		u.Error(err)
		return
	}

	// Remove the calendar from the upstream source
	err = event.GetCalendar().DeleteEvent(event, u.Tx.Queries())
	if err != nil {
		u.Error(err)
		return
	}

	// Delete event entry from the database
	err = u.Tx.Queries().DeleteEvent(userId, eventId)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

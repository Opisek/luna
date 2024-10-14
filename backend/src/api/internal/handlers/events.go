package handlers

import (
	"luna-backend/api/internal/context"
	"luna-backend/api/internal/util"
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
	//Settings primitives.EventSettings `json:"settings"` // TODO: REMOVE FROM PRODUCTION, TESTING ONLY
}

func GetEvents(c *gin.Context) {
	// Get config
	config := context.GetConfig(c)
	userId := context.GetUserId(c)
	calendarId, err := context.GetId(c, "calendar")
	if err != nil {
		config.Logger.Errorf("could not get calendar id: %v", err)
		util.Error(c, util.ErrorMalformedID)
		return
	}

	tx := context.GetTransaction(c)
	defer tx.Rollback(config.Logger)

	// Get the requested calendar
	calendar, err := tx.Queries().GetCalendar(userId, calendarId)
	if err != nil {
		config.Logger.Errorf("could not get events: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	// Get the associated events
	startStr := c.Query("start")
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		config.Logger.Warnf("could not get events: could not parse start time: %v", err)
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailTime)
		return
	}
	endStr := c.Query("end")
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		config.Logger.Warnf("could not get events: could not parse end time: %v", err)
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailTime)
		return
	}
	if startTime.After(endTime) {
		config.Logger.Warn("start time is after end time")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailTime)
		return
	}
	if endTime.Sub(startTime) > time.Hour*24*365 {
		endTime = startTime.Add(time.Hour * 24 * 365)
	}

	eventsFromCal, err := calendar.GetEvents(startTime, endTime)
	if err != nil {
		config.Logger.Errorf("could not get events: could not get events from calendar %v: %v", calendar.GetName(), err)
		util.Error(c, util.ErrorUnknown)
		return
	}

	events, err := tx.Queries().ReconcileEvents(eventsFromCal)
	if err != nil {
		config.Logger.Errorf("could not reconcile events: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	// Convert to exposed format
	convertedEvents := make([]exposedEvent, len(events))
	for i, event := range events {
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
			//Settings: event.GetSettings(),
		}
	}

	if tx.Commit(config.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": convertedEvents})
}

func GetEvent(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	eventId, err := context.GetId(c, "event")
	if err != nil {
		apiConfig.Logger.Errorf("could not get event id: %v", err)
		util.Error(c, util.ErrorMalformedID)
		return
	}

	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	// Get event
	event, err := tx.Queries().GetEvent(userId, eventId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get event: %v", err)
		util.Error(c, util.ErrorDatabase)
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

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusOK, convertedCal)
}

func PutEvent(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	calendarId, err := context.GetId(c, "calendar")
	if err != nil {
		apiConfig.Logger.Warn("missing or malformed calendar id")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailId)
		return
	}

	calendar, err := tx.Queries().GetCalendar(userId, calendarId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get calendar: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	eventName := c.PostForm("name")
	if eventName == "" {
		apiConfig.Logger.Warn("missing name")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailName)
		return
	}

	eventDesc := c.PostForm("desc")

	eventColor, err := types.ParseColor(c.PostForm("color"))
	if err != nil {
		apiConfig.Logger.Warnf("missing or malformed color: %v", err)
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailColor)
		return
	}

	eventDateAllDay := c.PostForm("date_all_day") == "true"

	eventDateStartStr := c.PostForm("date_start")
	eventDateStart, err := time.Parse(time.RFC3339, eventDateStartStr)
	if err != nil {
		apiConfig.Logger.Warnf("missing or malformed date start: %v", err)
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailDate)
		return
	}

	eventDateEndStr := c.PostForm("date_end")
	eventDateDurationStr := c.PostForm("date_duration")

	eventDateEnd, endErr := time.Parse(time.RFC3339, eventDateEndStr)
	eventDateDuration, durationErr := time.ParseDuration(eventDateDurationStr)

	if (endErr != nil && durationErr != nil) || (endErr == nil && durationErr == nil) {
		apiConfig.Logger.Warn("missing or malformed date start")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailDate)
		return
	}

	var date *types.EventDate
	if endErr == nil {
		date = types.NewEventDateFromEndTime(&eventDateStart, &eventDateEnd, eventDateAllDay, nil)
	} else {
		date = types.NewEventDateFromDuration(&eventDateStart, &eventDateDuration, eventDateAllDay, nil)
	}

	event, err := calendar.AddEvent(eventName, eventDesc, eventColor, date)
	if err != nil {
		apiConfig.Logger.Errorf("could not add event: %v", err)
		util.Error(c, util.ErrorInternal)
		return
	}

	err = tx.Queries().InsertEvent(event)
	if err != nil {
		apiConfig.Logger.Errorf("could not insert event: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": event.GetId().String()})
}

func PatchEvent(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	eventId, err := context.GetId(c, "event")
	if err != nil {
		apiConfig.Logger.Warn("missing or malformed event id")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailId)
		return
	}

	event, err := tx.Queries().GetEvent(userId, eventId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get event: %v", err)
		util.Error(c, util.ErrorDatabase)
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
		apiConfig.Logger.Warn("no values to change")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailFields)
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
			apiConfig.Logger.Warn("cannot specify both end and duration")
			util.ErrorDetailed(c, util.ErrorPayload, util.DetailDate)
			return
		} else if endErr == nil {
			newEventDate = types.NewEventDateFromEndTime(&eventDateStart, &eventDateEnd, eventDateAllDay, nil)
		} else {
			newEventDate = types.NewEventDateFromDuration(&eventDateStart, &eventDateDuration, eventDateAllDay, nil)
		}
	}

	newEvent, err := event.GetCalendar().EditEvent(event, newEventName, newEventDesc, newEventColor, newEventDate)
	if err != nil {
		apiConfig.Logger.Errorf("could not edit event: %v", err)
		util.Error(c, util.ErrorInternal)
		return
	}

	err = tx.Queries().UpdateEvent(newEvent)
	if err != nil {
		apiConfig.Logger.Errorf("could not update event: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	util.Success(c)
}

func DeleteEvent(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	eventId, err := context.GetId(c, "event")
	if err != nil {
		apiConfig.Logger.Errorf("could not get event id: %v", err)
		util.Error(c, util.ErrorMalformedID)
		return
	}

	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	// Get event first
	event, err := tx.Queries().GetEvent(userId, eventId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get event: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	// Remove the calendar from the upstream source
	err = event.GetCalendar().DeleteEvent(event)
	if err != nil {
		apiConfig.Logger.Errorf("could not delete event from remote source: %v", err)
		util.Error(c, util.ErrorInternal)
		return
	}

	// Delete event entry from the database
	err = tx.Queries().DeleteEvent(userId, eventId)
	if err != nil {
		apiConfig.Logger.Errorf("could not delete event: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	util.Success(c)
}

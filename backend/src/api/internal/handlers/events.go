package handlers

import (
	"errors"
	"fmt"
	"luna-backend/api/internal/config"
	"luna-backend/api/internal/context"
	"luna-backend/api/internal/util"
	"luna-backend/db"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"net/http"
	"sync"
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

func getEvents(config *config.Api, tx *db.Transaction, cals []primitives.Calendar, start time.Time, end time.Time) ([]primitives.Event, error) {
	// For each calendar, get its events
	events := make([][]primitives.Event, len(cals))
	errored := false

	waitGroup := sync.WaitGroup{}
	for i, cal := range cals {
		waitGroup.Add(1)
		go func(i int, cal primitives.Calendar) {
			defer waitGroup.Done()

			eventsFromCal, err := cal.GetEvents(start, end)
			if err != nil {
				errored = true
				config.Logger.Errorf("could not get events: could not get events from calendar %v: %v", cal.GetName(), err)
				return
			}

			events[i] = eventsFromCal
		}(i, cal)
	}

	// Combine (flatten) all calendars
	waitGroup.Wait()
	if errored {
		return nil, errors.New("at least one calendar failed to provide events")
	}

	combinedEvents := []primitives.Event{}
	for _, eventsFromCal := range events {
		combinedEvents = append(combinedEvents, eventsFromCal...)
	}

	// Reconcile with database
	combinedEvents, err := tx.Queries().ReconcileEvents(cals, combinedEvents)
	if err != nil {
		return nil, fmt.Errorf("could not reconcile events: %v", err)
	}

	return combinedEvents, nil
}

func GetEvents(c *gin.Context) {
	// Get config
	config := context.GetConfig(c)
	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(config.Logger)

	// Get all of user's sources
	srcs, err := getSources(config, tx, userId)
	if err != nil {
		config.Logger.Errorf("could not get events: %v", err)
		util.Error(c, util.ErrorUnknown)
		return
	}

	// Get their associated calendars
	cals, err := getCalendars(config, tx, srcs)
	if err != nil {
		config.Logger.Errorf("could not get events: %v", err)
		util.Error(c, util.ErrorUnknown)
		return
	}

	// Get their associated events
	// TODO: get time from the api request
	startTime, err := time.Parse(time.RFC3339, "2024-01-01T00:00:00+00:00")
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(time.RFC3339, "2025-01-01T00:00:00+00:00")
	if err != nil {
		panic(err)
	}

	events, err := getEvents(config, tx, cals, startTime, endTime)
	if err != nil {
		config.Logger.Errorf("could not get events: %v", err)
		util.Error(c, util.ErrorUnknown)
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

	c.JSON(http.StatusOK, convertedEvents)
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

	calendarIdStr := c.PostForm("calendar")
	calendarId, err := types.IdFromString(calendarIdStr)
	if err != nil {
		apiConfig.Logger.Error("missing or malformed calendar id")
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
		apiConfig.Logger.Error("missing name")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailName)
		return
	}

	eventDesc := c.PostForm("desc")

	eventColor, err := types.ParseColor(c.PostForm("color"))
	if err != nil {
		apiConfig.Logger.Errorf("missing or malformed color: %v", err)
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailColor)
		return
	}

	eventDateStartStr := c.PostForm("date_start")
	eventDateStart, err := time.Parse(time.RFC3339, eventDateStartStr)
	if err != nil {
		apiConfig.Logger.Errorf("missing or malformed date start: %v", err)
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailDate)
		return
	}

	eventDateEndStr := c.PostForm("date_end")
	eventDateDurationStr := c.PostForm("date_duration")

	eventDateEnd, endErr := time.Parse(time.RFC3339, eventDateEndStr)
	eventDateDuration, durationErr := time.ParseDuration(eventDateDurationStr)

	if (endErr != nil && durationErr != nil) || (endErr == nil && durationErr == nil) {
		apiConfig.Logger.Error("missing or malformed date start")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailDate)
		return
	}

	var date *types.EventDate
	if endErr == nil {
		date = types.NewEventDateFromEndTime(&eventDateStart, &eventDateEnd, nil)
	} else {
		date = types.NewEventDateFromDuration(&eventDateStart, &eventDateDuration, nil)
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
		apiConfig.Logger.Error("missing or malformed event id")
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

	eventDateStartStr := c.PostForm("date_start")
	eventDateStart, startErr := time.Parse(time.RFC3339, eventDateStartStr)

	eventDateEndStr := c.PostForm("date_end")
	eventDateDurationStr := c.PostForm("date_duration")

	eventDateEnd, endErr := time.Parse(time.RFC3339, eventDateEndStr)
	eventDateDuration, durationErr := time.ParseDuration(eventDateDurationStr)

	if newEventName == "" && newEventDesc == event.GetDesc() && colErr != nil && startErr != nil && endErr != nil && durationErr != nil {
		apiConfig.Logger.Error("no values to change")
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
		if startErr == nil {
			eventDateStart = *event.GetDate().Start()
		}
		if endErr != nil && durationErr == nil {
			if event.GetDate().SpecifyDuration() {
				eventDateDuration = *event.GetDate().Duration()
				newEventDate = types.NewEventDateFromDuration(&eventDateStart, &eventDateDuration, nil)
			} else {
				eventDateEnd = *event.GetDate().End()
				newEventDate = types.NewEventDateFromEndTime(&eventDateStart, &eventDateEnd, nil)
			}
		} else if endErr == nil && durationErr == nil {
			apiConfig.Logger.Error("cannot specify both end and duration")
			util.ErrorDetailed(c, util.ErrorPayload, util.DetailDate)
			return
		} else if endErr == nil {
			newEventDate = types.NewEventDateFromEndTime(&eventDateStart, &eventDateEnd, nil)
		} else {
			newEventDate = types.NewEventDateFromDuration(&eventDateStart, &eventDateDuration, nil)
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

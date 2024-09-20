package handlers

import (
	"errors"
	"luna-backend/api/internal/config"
	"luna-backend/api/internal/context"
	"luna-backend/db"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type exposedEvent struct {
	Id       types.ID                 `json:"id"`
	Calendar types.ID                 `json:"calendar"`
	Name     string                   `json:"name"`
	Desc     string                   `json:"desc"`
	Color    *types.Color             `json:"color"`
	Date     *types.EventDate         `json:"date"`
	Settings primitives.EventSettings `json:"settings"` // TODO: REMOVE FROM PRODUCTION, TESTING ONLY
}

func getEvents(config *config.Api, _ *db.Transaction, cals []primitives.Calendar, start time.Time, end time.Time) ([]primitives.Event, error) {
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
				config.Logger.Errorf("could not get events: could not get events from calendar %v: %v", cal.GetId().String(), err)
				return
			}

			events[i] = eventsFromCal
		}(i, cal)
	}

	// Combine (flatten) all calendars
	waitGroup.Wait()
	if errored {
		return nil, errors.New("at least one calendar failed to load")
	}

	combinedEvents := []primitives.Event{}
	for _, eventsFromCal := range events {
		combinedEvents = append(combinedEvents, eventsFromCal...)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get sources"})
		return
	}

	// Get their associated calendars
	cals, err := getCalendars(config, tx, srcs)
	if err != nil {
		config.Logger.Errorf("could not get events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get calendars"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get events"})
		return
	}

	// Reconcile with database
	events, err = tx.Queries().ReconcileEvents(cals, events)
	if err != nil {
		config.Logger.Errorf("could not reconcile events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reconcile events"})
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
			Settings: event.GetSettings(),
		}
	}

	if tx.Commit(config.Logger) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, convertedEvents)
}

func GetEvent(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	eventId, err := context.GetId(c, "event")
	if err != nil {
		apiConfig.Logger.Errorf("could not get event id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed or missing event id"})
		return
	}

	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	// Get event
	event, err := tx.Queries().GetEvent(userId, eventId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get event"})
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
		Settings: event.GetSettings(),
	}

	if tx.Commit(apiConfig.Logger) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, convertedCal)
}

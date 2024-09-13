package handlers

import (
	"errors"
	"luna-backend/api/internal/config"
	"luna-backend/api/internal/context"
	"luna-backend/interface/primitives"
	"luna-backend/types"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type exposedEvent struct {
	Id       types.ID     `json:"id"`
	Calendar types.ID     `json:"calendar"`
	Name     string       `json:"name"`
	Desc     string       `json:"desc"`
	Color    *types.Color `json:"color"`
}

func getEvents(config *config.Api, cals []primitives.Calendar) ([]primitives.Event, error) {
	// For each calendar, get its events
	events := make([][]primitives.Event, len(cals))
	errored := false

	waitGroup := sync.WaitGroup{}
	for i, cal := range cals {
		waitGroup.Add(1)
		go func(i int, cal primitives.Calendar) {
			defer waitGroup.Done()

			eventsFromCal, err := cal.GetEvents(time.Now(), time.Now()) // TODO: proper time filtering
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

	// Get all of user's sources
	srcs, err := getSources(config, userId)
	if err != nil {
		config.Logger.Errorf("could not get events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get sources"})
		return
	}

	// Get their associated calendars
	cals, err := getCalendars(config, srcs)
	if err != nil {
		config.Logger.Errorf("could not get events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get calendars"})
		return
	}

	// Get their associated events
	events, err := getEvents(config, cals)
	if err != nil {
		config.Logger.Errorf("could not get events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get events"})
		return
	}

	// Convert to exposed format
	convertedEvents := make([]exposedEvent, len(events))
	for i, event := range events {
		convertedEvents[i] = exposedEvent{
			Id:       event.GetId(),
			Calendar: event.GetCalendar(),
			Name:     event.GetName(),
			Desc:     event.GetDesc(),
			Color:    event.GetColor(),
		}
	}

	c.JSON(http.StatusOK, convertedEvents)
}

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

	"github.com/gin-gonic/gin"
)

type exposedCalendar struct {
	Id       types.ID                    `json:"id"`
	Source   types.ID                    `json:"source"`
	Name     string                      `json:"name"`
	Desc     string                      `json:"desc"`
	Color    *types.Color                `json:"color"`
	Settings primitives.CalendarSettings `json:"settings"` // TODO: REMOVE FROM PRODUCTION, TESTING ONLY
}

func getCalendars(config *config.Api, tx *db.Transaction, srcs []primitives.Source) ([]primitives.Calendar, error) {
	// For each source, get its calendars
	cals := make([][]primitives.Calendar, len(srcs))
	errored := false

	waitGroup := sync.WaitGroup{}
	for i, src := range srcs {
		waitGroup.Add(1)
		go func(i int, source primitives.Source) {
			defer waitGroup.Done()

			calsFromSource, err := tx.GetCalendars(source)
			if err != nil {
				errored = true
				config.Logger.Errorf("could not get calendars: %v", err)
				return
			}

			cals[i] = calsFromSource
		}(i, src)
	}

	// Combine (flatten) all calendars
	waitGroup.Wait()
	if errored {
		return nil, errors.New("at least one calendar failed to load")
	}

	combinedCals := []primitives.Calendar{}
	for _, calsFromSource := range cals {
		combinedCals = append(combinedCals, calsFromSource...)
	}

	return combinedCals, nil
}

func GetCalendars(c *gin.Context) {
	// Get config
	config := context.GetConfig(c)
	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(config.Logger)

	// Get all of user's sources
	srcs, err := getSources(config, tx, userId)
	if err != nil {
		config.Logger.Errorf("could not get calendars: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get sources"})
		return
	}

	// Get their associated calendars
	cals, err := getCalendars(config, tx, srcs)
	if err != nil {
		config.Logger.Errorf("could not get calendars: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get calendars"})
		return
	}

	// Convert to exposed format
	convertedCals := make([]exposedCalendar, len(cals))
	for i, cal := range cals {
		convertedCals[i] = exposedCalendar{
			Id:       cal.GetId(),
			Source:   cal.GetSource().GetId(),
			Name:     cal.GetName(),
			Desc:     cal.GetDesc(),
			Color:    cal.GetColor(),
			Settings: cal.GetSettings(),
		}
	}

	if tx.Commit(config.Logger) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, convertedCals)
}

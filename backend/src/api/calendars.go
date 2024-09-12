package api

import (
	"luna-backend/interface/primitives/sources"
	"luna-backend/types"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type exposedCalendar struct {
	Id     types.ID     `json:"id"`
	Source types.ID     `json:"source"`
	Name   string       `json:"name"`
	Desc   string       `json:"desc"`
	Color  *types.Color `json:"color"`
}

func getCalendars(c *gin.Context) {
	// Get config
	config := getConfig(c)
	if config == nil {
		return
	}

	userId := getUserId(c)

	// Get all of user's sources
	srcs, err := config.db.GetSources(userId)
	if err != nil {
		config.logger.Errorf("could not get calendars: could not get user's sources: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user's sources"})
		return
	}

	// For each source, get its calendars
	cals := make([][]exposedCalendar, len(srcs))
	errored := false

	waitGroup := sync.WaitGroup{}
	for i, src := range srcs {
		waitGroup.Add(1)
		go func(i int, source sources.Source) {
			defer waitGroup.Done()

			calsFromSource, err := source.GetCalendars()
			if err != nil {
				errored = true
				config.logger.Errorf("could not get calendars: could not get calendars from source: %v", err)
				return
			}

			convertedCals := make([]exposedCalendar, len(calsFromSource))
			for j, cal := range calsFromSource {
				convertedCals[j] = exposedCalendar{
					Id:     cal.GetId(),
					Source: source.GetId(),
					Name:   cal.GetName(),
					Desc:   cal.GetDesc(),
					Color:  cal.GetColor(),
				}
			}

			cals[i] = convertedCals
		}(i, src)
	}

	// Combine (flatten) all calendars
	waitGroup.Wait()
	if errored {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get calendars"})
		return
	}

	combinedCals := []exposedCalendar{}
	for _, calsFromSource := range cals {
		combinedCals = append(combinedCals, calsFromSource...)
	}

	c.JSON(http.StatusOK, combinedCals)
}

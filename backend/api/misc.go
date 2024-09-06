package api

import (
	"encoding/base64"
	"fmt"
	"luna-backend/sources"
	"luna-backend/types"
	"luna-backend/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func getCalendars(c *gin.Context) {
	calendars := make([]*types.Calendar, 0)

	for _, source := range util.Sources {
		cals, err := source.GetCalendars()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		calendars = append(calendars, cals...)
	}

	for _, calendar := range calendars {
		calendar.Id = base64.StdEncoding.EncodeToString([]byte(calendar.Source + calendar.Id))
	}

	c.JSON(http.StatusOK, calendars)
}

func getEvents(c *gin.Context) {
	// Get IDs
	combinedId, err := base64.StdEncoding.DecodeString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sourceId, err := sources.SourceIdFromBytes(combinedId[:36])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	source := util.GetSource(sourceId)

	calendarId := string(combinedId[36:])

	// Get Times

	startTime, err := time.Parse(time.RFC3339, c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time"})
		return
	}

	if !startTime.Before(endTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start time must be before end time"})
		return
	}
	if endTime.Sub(startTime) > 365*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "time range must be less than 1 year"})
		return
	}

	// Process request

	source.GetEvents(calendarId, startTime, endTime)

	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func getVersion(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"version": apiConfig.commonConfig.Version.String()})
}

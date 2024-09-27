package handlers

import (
	"fmt"
	"luna-backend/api/internal/config"
	"luna-backend/api/internal/context"
	"luna-backend/api/internal/util"
	"luna-backend/db"
	"luna-backend/interface/primitives"
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
	//Settings primitives.CalendarSettings `json:"settings"` // TODO: REMOVE FROM PRODUCTION, TESTING ONLY
}

func getCalendars(config *config.Api, tx *db.Transaction, srcs []primitives.Source) ([]primitives.Calendar, []bool, error) {
	// For each source, get its calendars
	cals := make([][]primitives.Calendar, len(srcs))
	success := make([]bool, len(srcs))

	waitGroup := sync.WaitGroup{}
	for i, src := range srcs {
		waitGroup.Add(1)
		go func(i int, source primitives.Source) {
			defer waitGroup.Done()

			calsFromSource, err := source.GetCalendars()
			success[i] = err == nil

			if err != nil {
				cals[i] = []primitives.Calendar{}
				config.Logger.Errorf("could not fetch calendars from source %v: %v", source.GetName(), err)
				return
			}

			cals[i] = calsFromSource
		}(i, src)
	}

	// Combine (flatten) all calendars
	waitGroup.Wait()

	combinedCals := []primitives.Calendar{}
	for _, calsFromSource := range cals {
		combinedCals = append(combinedCals, calsFromSource...)
	}

	// Reconcile with database
	combinedCals, err := tx.Queries().ReconcileCalendars(srcs, combinedCals)
	if err != nil {
		return nil, nil, fmt.Errorf("could not reconcile calendars: %v", err)
	}

	return combinedCals, success, nil
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
		util.Error(c, util.ErrorUnknown)
		return
	}

	// Get their associated calendars
	cals, succeeded, err := getCalendars(config, tx, srcs)
	if err != nil {
		config.Logger.Errorf("could not get calendars: %v", err)
		util.Error(c, util.ErrorUnknown)
		return
	}

	// Convert to exposed format
	convertedCals := make([]exposedCalendar, len(cals))
	for i, cal := range cals {
		convertedCals[i] = exposedCalendar{
			Id:     cal.GetId(),
			Source: cal.GetSource().GetId(),
			Name:   cal.GetName(),
			Desc:   cal.GetDesc(),
			Color:  cal.GetColor(),
			//Settings: cal.GetSettings(),
		}
	}

	if tx.Commit(config.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	errCnt := 0
	for _, success := range succeeded {
		if !success {
			errCnt++
		}
	}
	errIDs := make([]types.ID, errCnt)
	j := 0
	for i, success := range succeeded {
		if !success {
			errIDs[j] = srcs[i].GetId()
			j += 1
		}
	}

	c.JSON(http.StatusOK, gin.H{"calendars": convertedCals, "errored": errIDs})
}

func GetCalendar(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	calendarId, err := context.GetId(c, "calendar")
	if err != nil {
		apiConfig.Logger.Errorf("could not get calendar id: %v", err)
		util.Error(c, util.ErrorMalformedID)
		return
	}

	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	// Get calendar
	cal, err := tx.Queries().GetCalendar(userId, calendarId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get calendar: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	// Convert to exposed format
	convertedCal := exposedCalendar{
		Id:     cal.GetId(),
		Source: cal.GetSource().GetId(),
		Name:   cal.GetName(),
		Desc:   cal.GetDesc(),
		Color:  cal.GetColor(),
		//Settings: cal.GetSettings(),
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusOK, convertedCal)
}

func PutCalendar(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	sourceIdStr := c.PostForm("source")
	sourceId, err := types.IdFromString(sourceIdStr)
	if err != nil {
		apiConfig.Logger.Error("missing or malformed source id")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailId)
		return
	}

	source, err := tx.Queries().GetSource(userId, sourceId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get source: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	calName := c.PostForm("name")
	if calName == "" {
		apiConfig.Logger.Error("missing calendar name")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailName)
		return
	}

	calColor, err := types.ParseColor(c.PostForm("color"))
	if err != nil {
		apiConfig.Logger.Error("missing or malformed color")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailColor)
		return
	}

	cal, err := source.AddCalendar(calName, calColor)
	if err != nil {
		apiConfig.Logger.Errorf("could not add calendar: %v", err)
		util.Error(c, util.ErrorUnknown)
		return
	}

	err = tx.Queries().InsertCalendar(cal)
	if err != nil {
		apiConfig.Logger.Errorf("could not insert calendar: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": cal.GetId().String()})
}

func PatchCalendar(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	calendarId, err := context.GetId(c, "calendar")
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

	newCalName := c.PostForm("name")

	newCalColor, colErr := types.ParseColor(c.PostForm("color"))

	if newCalName == "" && colErr != nil {
		apiConfig.Logger.Error("no values to change")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailFields)
		return
	}

	if newCalName == "" {
		newCalName = calendar.GetName()
	}

	if colErr != nil {
		newCalColor = calendar.GetColor()
	}

	newCal, err := calendar.GetSource().EditCalendar(calendar, newCalName, newCalColor)
	if err != nil {
		apiConfig.Logger.Errorf("could not edit calendar: %v", err)
		util.Error(c, util.ErrorUnknown)
		return
	}

	err = tx.Queries().UpdateCalendar(newCal)
	if err != nil {
		apiConfig.Logger.Errorf("could not update calendar: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	util.Success(c)
}

func DeleteCalendar(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	calendarId, err := context.GetId(c, "calendar")
	if err != nil {
		apiConfig.Logger.Errorf("could not get calendar id: %v", err)
		util.Error(c, util.ErrorMalformedID)
		return
	}

	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	calendar, err := tx.Queries().GetCalendar(userId, calendarId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get calendar: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	err = calendar.GetSource().DeleteCalendar(calendar)
	if err != nil {
		apiConfig.Logger.Errorf("could not delete calendar: %v", err)
		util.Error(c, util.ErrorUnknown)
		return
	}

	err = tx.Queries().DeleteCalendar(userId, calendarId)
	if err != nil {
		apiConfig.Logger.Errorf("could not delete calendar: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	util.Success(c)
}

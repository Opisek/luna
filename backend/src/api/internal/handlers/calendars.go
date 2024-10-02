package handlers

import (
	"luna-backend/api/internal/context"
	"luna-backend/api/internal/util"
	"luna-backend/types"
	"net/http"

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

func GetCalendars(c *gin.Context) {
	// Get config
	config := context.GetConfig(c)
	userId := context.GetUserId(c)
	sourceId, err := context.GetId(c, "source")
	if err != nil {
		config.Logger.Errorf("could not get source id: %v", err)
		util.Error(c, util.ErrorMalformedID)
		return
	}

	tx := context.GetTransaction(c)
	defer tx.Rollback(config.Logger)

	// Get the specified source
	source, err := tx.Queries().GetSource(userId, sourceId)
	if err != nil {
		config.Logger.Errorf("could not get calendars: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	// Get the associated calendars
	calsFromSource, err := source.GetCalendars()

	if err != nil {
		config.Logger.Errorf("could not fetch calendars from source %v: %v", source.GetName(), err)
		util.Error(c, util.ErrorUnknown)
		return
	}

	cals, err := tx.Queries().ReconcileCalendars(calsFromSource)
	if err != nil {
		config.Logger.Errorf("could not reconcile calendars: %v", err)
		util.Error(c, util.ErrorDatabase)
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

	c.JSON(http.StatusOK, gin.H{"calendars": convertedCals})
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

	sourceId, err := context.GetId(c, "source")
	if err != nil {
		apiConfig.Logger.Warn("missing or malformed source id")
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
		apiConfig.Logger.Warn("missing calendar name")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailName)
		return
	}

	calColor, err := types.ParseColor(c.PostForm("color"))
	if err != nil {
		apiConfig.Logger.Warn("missing or malformed color")
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

	newCalName := c.PostForm("name")

	newCalColor, colErr := types.ParseColor(c.PostForm("color"))

	if newCalName == "" && colErr != nil {
		apiConfig.Logger.Warn("no values to change")
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

package api

import (
	"luna-backend/caldav"
	"luna-backend/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getCalendars(c *gin.Context) {
	calendars, err := caldav.GetCalendars(util.CaldavSettingsInstance)

	if err != nil {
		// TODO: debug levels, we don't want to expose the error message to the user
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, calendars)
	}
}

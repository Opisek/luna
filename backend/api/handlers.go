package api

import (
	"fmt"
	"luna-backend/types"
	"luna-backend/util"
	"net/http"

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

	c.JSON(http.StatusOK, calendars)
}

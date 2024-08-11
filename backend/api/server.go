package api

import "github.com/gin-gonic/gin"

func Run() {
	router := gin.Default()

	endpoints := router.Group("/api")
	endpoints.GET("/calendars", getCalendars)
	endpoints.GET("/events/:id", getEvents)

	router.Run(":3000")
}

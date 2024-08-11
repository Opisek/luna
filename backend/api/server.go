package api

import (
	"luna-backend/auth"

	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()

	endpoints := router.Group("/api")

	endpoints.POST("/login", auth.Login)

	authenticatedEndpoints := endpoints.Group("", auth.AuthMiddleware())

	sourcesEndpoints := authenticatedEndpoints.Group("/sources")
	sourcesEndpoints.GET("", getSources)
	sourcesEndpoints.PUT("", notImplemented)
	sourcesEndpoints.PATCH("/:sourceId", notImplemented)
	sourcesEndpoints.DELETE("/:sourceId", notImplemented)

	calendarsEndpoints := sourcesEndpoints.Group("/:sourceId/calendars")
	calendarsEndpoints.GET("", notImplemented)
	calendarsEndpoints.PUT("", notImplemented)
	calendarsEndpoints.PATCH("/:calendarId", notImplemented)
	calendarsEndpoints.DELETE("/:calendarId", notImplemented)

	eventEndpoints := calendarsEndpoints.Group("/:calendarId/events")
	eventEndpoints.GET("", notImplemented)
	eventEndpoints.PUT("", notImplemented)
	eventEndpoints.PATCH("/:eventId", notImplemented)
	eventEndpoints.DELETE("/:eventId", notImplemented)

	router.Run(":3000")
}

package api

import (
	"fmt"
	"luna-backend/api/internal/config"
	"luna-backend/api/internal/handlers"
	"luna-backend/common"
	"luna-backend/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewApi(db *db.Database, commonConfig *common.CommonConfig, logger *logrus.Entry) *config.Api {
	return config.NewApi(db, commonConfig, logger, run)
}

func run(api *config.Api) {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("apiConfig", api)
		c.Next()
	})

	// /api/*
	rawEndpoints := router.Group("/api")
	endpoints := rawEndpoints.Group("", handlers.TransactionMiddleware())

	endpoints.POST("/login", handlers.Login)
	endpoints.POST("/register", handlers.Register)
	endpoints.GET("/version", handlers.GetVersion)

	authenticatedEndpoints := endpoints.Group("", handlers.AuthMiddleware())

	// /api/sources/*
	sourcesEndpoints := authenticatedEndpoints.Group("/sources")
	sourcesEndpoints.GET("", handlers.GetSources)
	sourcesEndpoints.GET("/:sourceId", handlers.GetSource)
	sourcesEndpoints.PUT("", handlers.PutSource)
	sourcesEndpoints.PATCH("/:sourceId", handlers.PatchSource)
	sourcesEndpoints.DELETE("/:sourceId", handlers.DeleteSource)
	sourcesEndpoints.GET("/:sourceId/calendars", handlers.GetCalendars)
	sourcesEndpoints.PUT("/:sourceId/calendars", handlers.PutCalendar)

	// /api/calendars/*
	calendarsEndpoints := authenticatedEndpoints.Group("/calendars")
	calendarsEndpoints.GET("/:calendarId", handlers.GetCalendar)
	calendarsEndpoints.PATCH("/:calendarId", handlers.PatchCalendar)
	calendarsEndpoints.DELETE("/:calendarId", handlers.DeleteCalendar)
	calendarsEndpoints.GET("/:calendarId/events", handlers.GetEvents)
	calendarsEndpoints.PUT("/:calendarId/events", handlers.PutEvent)

	// /api/events/*
	eventEndpoints := authenticatedEndpoints.Group("/events")
	eventEndpoints.GET("/:eventId", handlers.GetEvent)
	eventEndpoints.PATCH("/:eventId", handlers.PatchEvent)
	eventEndpoints.DELETE("/:eventId", handlers.DeleteEvent)

	router.Run(fmt.Sprintf(":%d", api.CommonConfig.Env.API_PORT))
}

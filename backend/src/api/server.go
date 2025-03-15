package api

import (
	"fmt"
	"luna-backend/api/internal/config"
	"luna-backend/api/internal/handlers"
	"luna-backend/common"
	"luna-backend/db"
	"time"

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
	rawEndpoints := router.Group("/api")

	// /api/* (with no transactions)
	noDatabaseEndpoints := rawEndpoints.Group("",
		handlers.ContextMiddleware(3*time.Second),
		gin.Recovery(),
	)

	noDatabaseEndpoints.GET("/version", handlers.GetVersion)

	// /api/* (long-running authentication)
	authenticationEndpoints := rawEndpoints.Group("",
		handlers.ContextMiddleware(20*time.Second),
		gin.Recovery(),
		handlers.TransactionMiddleware(),
	)

	authenticationEndpoints.POST("/login", handlers.Login)
	authenticationEndpoints.POST("/register", handlers.Register)

	// /api/* the rest
	endpoints := noDatabaseEndpoints.Group("", handlers.TransactionMiddleware())

	endpoints.GET("/health", handlers.GetHealth)

	// everything past here requires the user to be logged in
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

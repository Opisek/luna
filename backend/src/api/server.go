package api

import (
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
	endpoints := router.Group("/api")

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

	// /api/calendars/*
	calendarsEndpoints := authenticatedEndpoints.Group("/calendars")
	calendarsEndpoints.GET("", handlers.GetCalendars)
	calendarsEndpoints.GET("/:calendarId", handlers.NotImplemented)
	calendarsEndpoints.PUT("", handlers.NotImplemented)
	calendarsEndpoints.PATCH("/:calendarId", handlers.NotImplemented)
	calendarsEndpoints.DELETE("/:calendarId", handlers.NotImplemented)

	// /api/events/*
	eventEndpoints := authenticatedEndpoints.Group("/events")
	eventEndpoints.GET("", handlers.GetEvents)
	eventEndpoints.GET("/:eventId", handlers.NotImplemented)
	eventEndpoints.PUT("", handlers.NotImplemented)
	eventEndpoints.PATCH("/:eventId", handlers.NotImplemented)
	eventEndpoints.DELETE("/:eventId", handlers.NotImplemented)

	router.Run(":3000")
}

package api

import (
	"luna-backend/common"
	"luna-backend/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Api struct {
	db           *db.Database
	commonConfig *common.CommonConfig
	logger       *logrus.Entry
}

func NewApi(db *db.Database, commonConfig *common.CommonConfig, logger *logrus.Entry) *Api {
	return &Api{
		db:           db,
		commonConfig: commonConfig,
		logger:       logger,
	}
}

func (api *Api) Run() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("apiConfig", api)
		c.Next()
	})

	// /api/*
	endpoints := router.Group("/api")

	endpoints.POST("/login", login)
	endpoints.POST("/register", register)
	endpoints.GET("/version", getVersion)

	authenticatedEndpoints := endpoints.Group("", authMiddleware())

	// /api/sources/*
	sourcesEndpoints := authenticatedEndpoints.Group("/sources")
	sourcesEndpoints.GET("", getSources)
	sourcesEndpoints.GET("/:sourceId", getSource)
	sourcesEndpoints.PUT("", putSource)
	sourcesEndpoints.PATCH("/:sourceId", patchSource)
	sourcesEndpoints.DELETE("/:sourceId", deleteSource)

	// /api/calendars/*
	//calendarsEndpoints := sourcesEndpoints.Group("/:sourceId/calendars") // old approach, however now we have unique IDs
	calendarsEndpoints := authenticatedEndpoints.Group("/calendars")
	calendarsEndpoints.GET("", getCalendars)
	calendarsEndpoints.GET("/:calendarId", notImplemented)
	calendarsEndpoints.PUT("", notImplemented)
	calendarsEndpoints.PATCH("/:calendarId", notImplemented)
	calendarsEndpoints.DELETE("/:calendarId", notImplemented)

	// /api/sources/.../calendars/.../events/*
	eventEndpoints := calendarsEndpoints.Group("/:calendarId/events")
	eventEndpoints.GET("", notImplemented)
	eventEndpoints.PUT("", notImplemented)
	eventEndpoints.PATCH("/:eventId", notImplemented)
	eventEndpoints.DELETE("/:eventId", notImplemented)

	router.Run(":3000")
}

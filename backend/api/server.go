package api

import (
	"luna-backend/auth"
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

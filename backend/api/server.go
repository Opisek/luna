package api

import (
	"luna-backend/auth"
	"luna-backend/common"
	"luna-backend/db"
	"net/http"

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

	authenticatedEndpoints := endpoints.Group("", auth.AuthMiddleware())

	// /api/sources/*
	sourcesEndpoints := authenticatedEndpoints.Group("/sources")
	sourcesEndpoints.GET("", getSources)
	sourcesEndpoints.PUT("", notImplemented)
	sourcesEndpoints.PATCH("/:sourceId", notImplemented)
	sourcesEndpoints.DELETE("/:sourceId", notImplemented)

	// /api/sources/.../calendars/*
	calendarsEndpoints := sourcesEndpoints.Group("/:sourceId/calendars")
	calendarsEndpoints.GET("", notImplemented)
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

func getConfig(c *gin.Context) *Api {
	// TODO: consider changing to "MustGet"
	apiConfig, err := c.Get("apiConfig")
	if !err {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "context error"})
		return nil
	}
	return apiConfig.(*Api)
}

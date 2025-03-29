package api

import (
	"fmt"
	middleware "luna-backend/api/internal"
	"luna-backend/api/internal/handlers"
	"luna-backend/api/internal/util"
	"luna-backend/config"
	"luna-backend/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewApi(db *db.Database, commonConfig *config.CommonConfig, logger *logrus.Entry) *util.Api {
	return util.NewApi(db, commonConfig, logger, run)
}

func run(api *util.Api) {
	router := gin.Default()
	rawEndpoints := router.Group("/api")

	// /api/* (with no transactions)
	noDatabaseEndpoints := rawEndpoints.Group("", middleware.RequestSetup(3*time.Second, api.Db, false, api.CommonConfig, api.Logger))

	noDatabaseEndpoints.GET("/version", handlers.GetVersion)

	// /api/* (long-running authentication)
	authenticationEndpoints := rawEndpoints.Group("", middleware.RequestSetup(30*time.Second, api.Db, true, api.CommonConfig, api.Logger))

	authenticationEndpoints.POST("/login", handlers.Login)
	authenticationEndpoints.POST("/register", handlers.Register)

	// /api/* the rest
	endpoints := rawEndpoints.Group("", middleware.RequestSetup(3*time.Second, api.Db, true, api.CommonConfig, api.Logger))

	endpoints.GET("/health", handlers.GetHealth)

	// everything past here requires the user to be logged in
	noDatabaseAuthenticatedEndpoints := noDatabaseEndpoints.Group("", middleware.RequireAuth())
	authenticatedEndpoints := endpoints.Group("", middleware.RequireAuth())
	administratorEndpoints := authenticatedEndpoints.Group("", middleware.RequireAdmin())

	// /api/user
	userEndpoints := authenticatedEndpoints.Group("/user")
	userEndpoints.GET("", handlers.GetUserData)
	userEndpoints.PATCH("", handlers.PatchUserData)
	userEndpoints.DELETE("", handlers.DeleteUser)

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

	// /api/files/*
	fileEndpoints := authenticatedEndpoints.Group("/files")
	fileEndpoints.GET("/:fileId", handlers.GetFile)
	fileEndpoints.HEAD("/:fileId", handlers.GetFile)

	// /api/settings
	userSettingsEndpoints := authenticatedEndpoints.Group("/settings/user")
	userSettingsEndpoints.GET("", handlers.GetUserSettings)
	userSettingsEndpoints.GET("/:settingKey", handlers.GetUserSetting)
	userSettingsEndpoints.PATCH("/:settingKey", handlers.PatchUserSetting)
	userSettingsEndpoints.DELETE("/:settingKey", handlers.ResetUserSetting)

	globalSettingsEndpoints := administratorEndpoints.Group("/settings/global")
	globalSettingsEndpointsPublic := authenticatedEndpoints.Group("/settings/global")
	globalSettingsEndpointsPublic.GET("", handlers.GetGlobalSettings)
	globalSettingsEndpointsPublic.GET("/:settingKey", handlers.GetGlobalSetting)
	globalSettingsEndpoints.PATCH("/:settingKey", handlers.PatchGlobalSetting)
	globalSettingsEndpoints.DELETE("/:settingKey", handlers.ResetGlobalSetting)

	// /api/* the rest
	noDatabaseAuthenticatedEndpoints.POST("/url", handlers.CheckUrl)

	// Run the server
	router.Run(fmt.Sprintf(":%d", api.CommonConfig.Env.API_PORT))
}

package api

import (
	"fmt"
	middleware "luna-backend/api/internal"
	"luna-backend/api/internal/handlers"
	"luna-backend/api/internal/util"
	"luna-backend/config"
	"luna-backend/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewApi(db *db.Database, commonConfig *config.CommonConfig, logger *logrus.Entry) *util.Api {
	return util.NewApi(db, commonConfig, logger, run)
}

func run(api *util.Api) {
	if api.CommonConfig.Env.DEVELOPMENT {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1", "localhost", "::1", "192.168.0.0/16", "172.16.0.0/12", "172.17.0.0/16", "172.18.0.0/16", "10.0.0.8/8"})
	rawEndpoints := router.Group("/api")

	// /api/* (with no transactions)
	noDatabaseEndpoints := rawEndpoints.Group("", middleware.RequestSetup(api.CommonConfig.Env.REQUEST_TIMEOUT_DEFAULT, api.Db, false, api.CommonConfig, api.Logger))

	noDatabaseEndpoints.GET("/version", handlers.GetVersion)

	// /api/* (long-running authentication)
	longRunningEndpoints := rawEndpoints.Group("", middleware.RequestSetup(api.CommonConfig.Env.REQUEST_TIMEOUT_AUTHENTICATION, api.Db, true, api.CommonConfig, api.Logger))

	longRunningEndpoints.POST("/login", handlers.Login)
	longRunningEndpoints.POST("/register", handlers.Register)

	// /api/* the rest
	endpoints := rawEndpoints.Group("", middleware.RequestSetup(api.CommonConfig.Env.REQUEST_TIMEOUT_DEFAULT, api.Db, true, api.CommonConfig, api.Logger))

	endpoints.GET("/health", handlers.GetHealth)

	// everything past here requires the user to be logged in
	noDatabaseAuthenticatedEndpoints := noDatabaseEndpoints.Group("", middleware.RequireAuth())
	authenticatedEndpoints := endpoints.Group("", middleware.RequireAuth())
	longRunningAuthenticatedEndpoints := longRunningEndpoints.Group("", middleware.RequireAuth())
	administratorEndpoints := authenticatedEndpoints.Group("", middleware.RequireAdmin())

	// /api/users
	userEndpoints := authenticatedEndpoints.Group("/users")
	longRunningUserEndpoints := longRunningAuthenticatedEndpoints.Group("/users")
	userEndpoints.GET("/:userId", handlers.GetUserData)
	longRunningUserEndpoints.PATCH("/:userId", handlers.PatchUserData)
	longRunningUserEndpoints.DELETE("/:userId", handlers.DeleteUser)

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
	userSettingsEndpoints := userEndpoints.Group("/:userId/settings")
	userSettingsEndpoints.GET("", handlers.GetUserSettings)
	userSettingsEndpoints.GET("/:settingKey", handlers.GetUserSetting)
	userSettingsEndpoints.PATCH("", handlers.PatchUserSettings)
	userSettingsEndpoints.DELETE("", handlers.ResetUserSettings)
	userSettingsEndpoints.DELETE("/:settingKey", handlers.ResetUserSetting)

	globalSettingsEndpoints := administratorEndpoints.Group("/settings")
	globalSettingsEndpointsPublic := authenticatedEndpoints.Group("/settings")
	globalSettingsEndpointsPublic.GET("", handlers.GetGlobalSettings)
	globalSettingsEndpointsPublic.GET("/:settingKey", handlers.GetGlobalSetting)
	globalSettingsEndpoints.PATCH("", handlers.PatchGlobalSettings)
	globalSettingsEndpoints.DELETE("", handlers.ResetGlobalSettings)
	globalSettingsEndpoints.DELETE("/:settingKey", handlers.ResetGlobalSetting)

	// /api/* the rest
	noDatabaseAuthenticatedEndpoints.POST("/url", handlers.CheckUrl)

	// Run the server
	router.Run(fmt.Sprintf(":%d", api.CommonConfig.Env.API_PORT))
}

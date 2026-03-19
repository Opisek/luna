package api

import (
	"fmt"
	middleware "luna-backend/api/internal"
	"luna-backend/api/internal/handlers"
	"luna-backend/api/internal/util"
	"luna-backend/config"
	"luna-backend/db"
	"luna-backend/types"

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
	noDatabaseEndpoints := rawEndpoints.Group("",
		middleware.RequestSetup(api.CommonConfig.Env.REQUEST_TIMEOUT_DEFAULT, api.Db, false, api.CommonConfig, api.Logger),
	)

	noDatabaseEndpoints.GET("/version", handlers.GetVersion)

	// /api/* (long-running authentication)
	authEndpoints := rawEndpoints.Group("",
		middleware.RequestSetup(api.CommonConfig.Env.REQUEST_TIMEOUT_AUTHENTICATION, api.Db, true, api.CommonConfig, api.Logger),
		middleware.DynamicThrottle(api.Throttle),
	)

	authEndpoints.POST("/login", handlers.Login)
	authEndpoints.POST("/register", handlers.Register)

	// /api/* the rest
	endpoints := rawEndpoints.Group("",
		middleware.RequestSetup(api.CommonConfig.Env.REQUEST_TIMEOUT_DEFAULT, api.Db, true, api.CommonConfig, api.Logger),
	)

	endpoints.GET("/health", handlers.GetHealth)
	endpoints.GET("/register/enabled", handlers.RegistrationEnabled)

	// everything past here requires the user to be logged in
	authenticatedEndpoints := endpoints.Group("", middleware.RequireAuth())
	longRunningAuthenticatedEndpoints := authEndpoints.Group("", middleware.RequireAuth())
	administratorEndpoints := authenticatedEndpoints.Group("", middleware.RequireAdmin())

	// /api/users/*
	userEndpoints := authenticatedEndpoints.Group("/users", middleware.RequirePermissions(types.PermManageUsers))
	longRunningUserEndpoints := longRunningAuthenticatedEndpoints.Group("/users") // user endpoints that require password verification
	administrativeUserEndpoints := administratorEndpoints.Group("/users")

	userEndpoints.GET("/:userId", handlers.GetUser)
	userEndpoints.GET("", handlers.GetUsers)
	administrativeUserEndpoints.POST("/:userId/enable", handlers.EnableUser)
	administrativeUserEndpoints.POST("/:userId/disable", handlers.DisableUser)
	longRunningUserEndpoints.PATCH("/:userId", handlers.PatchUserData)
	longRunningUserEndpoints.DELETE("/:userId", handlers.DeleteUser)

	// /api/users/settings/*
	userSettingsEndpoints := userEndpoints.Group("/:userId/settings", middleware.RequirePermissions(types.PermManageUserSettings))
	userSettingsEndpoints.GET("", handlers.GetUserSettings)
	userSettingsEndpoints.GET("/:settingKey", handlers.GetUserSetting)
	userSettingsEndpoints.PATCH("", handlers.PatchUserSettings)
	userSettingsEndpoints.DELETE("", handlers.ResetUserSettings)
	userSettingsEndpoints.DELETE("/:settingKey", handlers.ResetUserSetting)

	// /api/sources/*
	sourcesEndpoints := authenticatedEndpoints.Group("/sources")

	sourcesEndpoints.GET("", middleware.RequirePermissions(types.PermReadSources), handlers.GetSources)
	sourcesEndpoints.GET("/:sourceId", middleware.RequirePermissions(types.PermReadSources), handlers.GetSource)
	sourcesEndpoints.PUT("", middleware.RequirePermissions(types.PermAddSources), handlers.PutSource)
	sourcesEndpoints.PATCH("/:sourceId", middleware.RequirePermissions(types.PermEditSources), handlers.PatchSource)
	sourcesEndpoints.DELETE("/:sourceId", middleware.RequirePermissions(types.PermDeleteSources), handlers.DeleteSource)
	sourcesEndpoints.GET("/:sourceId/calendars", middleware.RequirePermissions(types.PermReadCalendars), handlers.GetCalendars)
	sourcesEndpoints.PUT("/:sourceId/calendars", middleware.RequirePermissions(types.PermAddCalendars), handlers.PutCalendar)
	sourcesEndpoints.POST("/:sourceId/order", middleware.RequirePermissions(types.PermEditSources), handlers.ChangeSourceDisplayOrder)

	// /api/calendars/*
	calendarsEndpoints := authenticatedEndpoints.Group("/calendars")
	calendarsEndpoints.GET("/:calendarId", middleware.RequirePermissions(types.PermReadCalendars), handlers.GetCalendar)
	calendarsEndpoints.PATCH("/:calendarId", middleware.RequirePermissions(types.PermEditCalendars), handlers.PatchCalendar)
	calendarsEndpoints.DELETE("/:calendarId", middleware.RequirePermissions(types.PermDeleteCalendars), handlers.DeleteCalendar)
	calendarsEndpoints.GET("/:calendarId/events", middleware.RequirePermissions(types.PermReadEvents), handlers.GetEvents)
	calendarsEndpoints.PUT("/:calendarId/events", middleware.RequirePermissions(types.PermAddEvents), handlers.PutEvent)

	// /api/events/*
	eventEndpoints := authenticatedEndpoints.Group("/events")
	eventEndpoints.GET("/:eventId", middleware.RequirePermissions(types.PermReadEvents), handlers.GetEvent)
	eventEndpoints.PATCH("/:eventId", middleware.RequirePermissions(types.PermEditEvents), handlers.PatchEvent)
	eventEndpoints.DELETE("/:eventId", middleware.RequirePermissions(types.PermDeleteEvents), handlers.DeleteEvent)

	// /api/files/*
	fileEndpoints := authenticatedEndpoints.Group("/files")
	fileEndpoints.GET("/:fileId", handlers.GetFile)
	fileEndpoints.HEAD("/:fileId", handlers.GetFile)

	// /api/settings/*
	globalSettingsEndpoints := administratorEndpoints.Group("/settings", middleware.RequirePermissions(types.PermManageGlobalSettings))
	globalSettingsEndpointsPublic := authenticatedEndpoints.Group("/settings")
	globalSettingsEndpointsPublic.GET("", handlers.GetGlobalSettings)
	globalSettingsEndpointsPublic.GET("/:settingKey", handlers.GetGlobalSetting)
	globalSettingsEndpoints.PATCH("", handlers.PatchGlobalSettings)
	globalSettingsEndpoints.DELETE("", handlers.ResetGlobalSettings)
	globalSettingsEndpoints.DELETE("/:settingKey", handlers.ResetGlobalSetting)

	// /api/sessions/*
	sessionEndpoints := authenticatedEndpoints.Group("/sessions")
	sessionEndpoints.GET("/valid", handlers.IsSessionValid)
	sessionEndpoints.GET("/:sessionId/permissions", handlers.GetSessionPermissions)

	sessionManagementEndpoints := sessionEndpoints.Group("", middleware.RequirePermissions(types.PermManageSessions))
	sessionManagementEndpoints.GET("", handlers.GetSessions)
	sessionManagementEndpoints.PUT("", handlers.PutSession)
	sessionManagementEndpoints.PATCH("/:sessionId", handlers.PatchSession)
	sessionManagementEndpoints.DELETE("/:sessionId", handlers.DeleteSession)
	sessionManagementEndpoints.DELETE("", handlers.DeleteSessions)

	// /api/invites/*
	inviteEndpoints := administratorEndpoints.Group("/invites", middleware.RequirePermissions(types.PermManageInvites))
	inviteEndpoints.GET("", handlers.GetInvites)
	inviteEndpoints.GET("/:inviteId/qr", handlers.GetInviteQrCode)
	inviteEndpoints.PUT("", handlers.PutInvite)
	inviteEndpoints.DELETE("/:inviteId", handlers.DeleteInvite)
	inviteEndpoints.DELETE("", handlers.DeleteInvites)

	// /api/oauth/*
	oauthEndpoints := authenticatedEndpoints.Group("/oauth", middleware.RequirePermissions(types.PermManageUsers))
	oauthAdminEndpoints := administratorEndpoints.Group("/oauth", middleware.RequirePermissions(types.PermManageOauthClients))

	// /api/oauth/clients/*
	oauthClientEndpoints := oauthEndpoints.Group("/clients")
	oauthClientAdminEndpoints := oauthAdminEndpoints.Group("/clients")

	oauthClientEndpoints.GET("", handlers.GetOauthClients) // users must be able to see what auth providers they may use, but we must not return client secrets here
	oauthClientAdminEndpoints.GET("/:clientId", handlers.GetOauthClient)
	oauthClientAdminEndpoints.PUT("", handlers.PutOauthClient)
	oauthClientAdminEndpoints.PATCH("/:clientId", handlers.PatchOauthClient)
	oauthClientAdminEndpoints.DELETE("/:clientId", handlers.DeleteOauthClient)

	// /api/oauth/authorization/*
	oauthAuthRequestsEndpoints := oauthEndpoints.Group("/authorization")

	oauthAuthRequestsEndpoints.PUT("/:clientId", handlers.CreateOauthAuthorizationRequest)
	oauthAuthRequestsEndpoints.POST("/:requestId", handlers.FinalizeOauthAuthorizationRequest)
	oauthAuthRequestsEndpoints.DELETE("/:requestId", handlers.CancelOauthAuthorizationRequest)

	// /api/oauth/tokens/*
	oauthTokensEndpoints := oauthEndpoints.Group("/tokens")

	oauthTokensEndpoints.GET("", handlers.GetOauthClientsWithTokens)

	// /api/* the rest
	authenticatedEndpoints.POST("/url", handlers.CheckUrl)

	// Run the server
	router.Run(fmt.Sprintf(":%d", api.CommonConfig.Env.API_PORT))
}

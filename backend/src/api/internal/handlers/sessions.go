package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/errors"
	"luna-backend/types"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSessions(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)
	sessionId := util.GetSessionId(c)

	sessions, tr := u.Tx.Queries().GetSessions(userId)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{
		"sessions": sessions,
		"current":  sessionId,
	})
}

func PutSession(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	// Parse the supplied password
	password := c.PostForm("password")
	if password == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing password"),
		)
		return
	}
	// Get the user's password
	savedPassword, err := u.Tx.Queries().GetPassword(userId)
	if err != nil {
		u.Error(err.Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "Could not get password for user %v", userId.String()).
			Append(errors.LvlPlain, "Invalid credentials"),
		)
		return
	}

	// Verify the password
	if !auth.VerifyPassword(password, savedPassword) {
		u.Error(errors.New().Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "Wrong password").
			Append(errors.LvlPlain, "Invalid credentials"),
		)
		return
	}

	// Create the API token
	apiTokenName := c.Request.FormValue("name")
	if apiTokenName == "" {
		u.Error(errors.New().Status(400).
			Append(errors.LvlPlain, "Name may not be empty"),
		)
		return
	}

	session := &types.Session{
		UserId:       userId,
		UserAgent:    apiTokenName,
		IpAddress:    net.ParseIP(c.ClientIP()),
		IsShortLived: false,
		IsApi:        true,
	}
	err = u.Tx.Queries().InsertSession(session)
	if err != nil {
		u.Error(err.
			Append(errors.LvlWordy, "Could not create API session").
			AltStr(errors.LvlBroad, "Could not create API key"),
		)
		return
	}

	// Generate the token
	token, err := auth.NewToken(u.Config, u.Tx, userId, session.SessionId)
	if err != nil {
		u.Error(err.
			Append(errors.LvlWordy, "Could not generate API token").
			AltStr(errors.LvlBroad, "Could not create API key"),
		)
	}

	u.Success(&gin.H{
		"token": token,
	})
}

func DeleteSession(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)
	currentSessionId := util.GetSessionId(c)

	var sessionId types.ID
	var tr *errors.ErrorTrace
	if c.Param("sessionId") == "current" {
		sessionId = currentSessionId
	} else {
		sessionId, tr = util.GetId(c, "session")
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	tr = u.Tx.Queries().DeleteSession(userId, sessionId)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(nil)
}

func DeleteSessions(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	tr := u.Tx.Queries().DeleteSessions(userId)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(nil)
}

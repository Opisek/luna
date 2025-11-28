package handlers

import (
	"encoding/json"
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/perms"
	"luna-backend/types"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsSessionValid(c *gin.Context) {
	// All the validation is done by the various middleware instances
	util.GetUtil(c).Success(&gin.H{
		"valid": true,
	})
}

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

func GetSessionPermissions(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)
	isAdmin, tr := u.Tx.Queries().IsAdmin(userId)
	if tr != nil {
		u.Error(tr)
		return
	}

	var permissions *perms.TokenPermissions

	if c.Param("sessionId") == "current" {
		permissions = util.GetPermissions(c)
	} else if util.HasPermission(c, perms.ManageSessions) {
		sessionId, tr := util.GetId(c, "session")
		if tr != nil {
			u.Error(tr)
			return
		}

		sessionPermissions, tx := u.Tx.Queries().GetTokenPermissions(sessionId)
		if tx != nil {
			u.Error(tx)
			return
		}
		permissions = sessionPermissions
	} else {
		u.Error(errors.New().Status(http.StatusForbidden).
			Append(errors.LvlPlain, "You are not authorized to perform this action").
			AltStr(errors.LvlWordy, "You are missing one or more permissions").
			Append(errors.LvlDebug, "Missing permission: %v", perms.ManageSessions),
		)
		return
	}

	u.Success(&gin.H{
		"is_admin":    isAdmin,
		"permissions": permissions.ToList(),
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
	savedPassword, tr := u.Tx.Queries().GetPassword(userId)
	if tr != nil {
		u.Error(tr.Status(http.StatusUnauthorized).
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

	// Create the API session
	apiTokenName := c.Request.FormValue("name")
	if apiTokenName == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Name may not be empty"),
		)
		return
	}

	requestedPerms := c.Request.FormValue("permissions")
	var parsedPerms []string
	// Parse list of permissions []:
	err := json.Unmarshal([]byte(requestedPerms), &parsedPerms)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing or malformed permissions list"),
		)
		return
	}

	secret, tr := crypto.GenerateRandomBytes(256)
	if tr != nil {
		u.Error(tr.
			Append(errors.LvlWordy, "Could not generate random bytes").
			AltStr(errors.LvlBroad, "Could not create API key"),
		)
		return
	}

	session := &types.Session{
		UserId:           userId,
		UserAgent:        apiTokenName,
		InitialIpAddress: net.ParseIP(c.ClientIP()),
		LastIpAddress:    net.ParseIP(c.ClientIP()),
		IsShortLived:     false,
		IsApi:            true,
		SecretHash:       crypto.GetSha256Hash(secret),
		Permissions:      perms.FromStringList(parsedPerms),
	}
	tr = u.Tx.Queries().InsertSession(session)
	if tr != nil {
		u.Error(tr.
			Append(errors.LvlWordy, "Could not create API session").
			AltStr(errors.LvlBroad, "Could not create API key"),
		)
		return
	}

	tr = u.Tx.Queries().UpdateTokenPermissions(session.SessionId, session.Permissions)
	if tr != nil {
		u.Error(tr.
			Append(errors.LvlWordy, "Could not set API token permissions").
			AltStr(errors.LvlBroad, "Could not create API key"),
		)
		return
	}

	// Generate the token
	token, tr := auth.NewToken(u.Config, u.Tx, userId, session.SessionId, secret)
	if tr != nil {
		u.Error(tr.
			Append(errors.LvlWordy, "Could not generate API token").
			AltStr(errors.LvlBroad, "Could not create API key"),
		)
	}

	u.Success(&gin.H{
		"token": token,
	})
}

func PatchSession(c *gin.Context) {
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
	savedPassword, tr := u.Tx.Queries().GetPassword(userId)
	if tr != nil {
		u.Error(tr.Status(http.StatusUnauthorized).
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

	currentSessionId := util.GetSessionId(c)

	var sessionId types.ID
	if c.Param("sessionId") == "current" {
		sessionId = currentSessionId
	} else {
		sessionId, tr = util.GetId(c, "session")
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	session, tr := u.Tx.Queries().GetSession(userId, sessionId)
	if tr != nil {
		u.Error(tr)
		return
	}

	if !session.IsApi {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Cannot modify user sessions"),
		)
		return
	}

	apiTokenName := c.Request.FormValue("name")
	if apiTokenName == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Name may not be empty"),
		)
		return
	}

	differentName := apiTokenName != session.UserAgent
	if differentName {
		tr = u.Tx.Queries().UpdateSession(session)
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	session.UserAgent = apiTokenName

	requestedPerms := c.Request.FormValue("permissions")
	var parsedPerms []string
	// Parse list of permissions []:
	err := json.Unmarshal([]byte(requestedPerms), &parsedPerms)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing or malformed permissions list"),
		)
		return
	}
	session.Permissions = perms.FromStringList(parsedPerms)

	currentPermissions, tr := u.Tx.Queries().GetTokenPermissions(session.SessionId)
	if tr != nil {
		u.Error(tr)
		return
	}

	differentPermissions := !session.Permissions.Equals(currentPermissions)

	if differentPermissions {
		tr = u.Tx.Queries().UpdateTokenPermissions(session.SessionId, session.Permissions)
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	if !differentName && !differentPermissions {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Nothing to change"),
		)
		return
	}

	u.Success(nil)
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

	var tr *errors.ErrorTrace
	switch c.Query("type") {
	case "all":
		tr = u.Tx.Queries().DeleteSessions(userId)
	case "api":
		tr = u.Tx.Queries().DeleteApiSessions(userId)
	case "user":
		fallthrough
	default:
		tr = u.Tx.Queries().DeleteUserSessions(userId)
	}

	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(nil)
}

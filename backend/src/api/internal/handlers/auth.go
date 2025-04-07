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

// Error messages are intentionally kept vague in lower verbosity levels,
// because detailed error messages about authenticatino checks might pose a
// security risk.

func Login(c *gin.Context) {
	// Parsing
	u := util.GetUtil(c)

	credentials := auth.BasicAuth{}
	if err := c.ShouldBind(&credentials); err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse credentials").
			Append(errors.LvlWordy, "Malformed request").
			Append(errors.LvlBroad, "Could not log in"),
		)
		return
	}

	usernameErr := util.IsValidUsername(credentials.Username)
	passwordErr := util.IsValidPassword(credentials.Password)
	if usernameErr != nil || passwordErr != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, usernameErr).AndErr(passwordErr).
			Append(errors.LvlDebug, "Input did not pass validation").
			Append(errors.LvlWordy, "Malformed request").
			Append(errors.LvlBroad, "Could not log in"),
		)
		return
	}

	// Check if the user exists
	userId, err := u.Tx.Queries().GetUserIdFromUsername(credentials.Username)
	if err != nil {
		u.Error(err.Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "Could not find ID for user %v", credentials.Username).
			Append(errors.LvlPlain, "Invalid credentials").
			Append(errors.LvlBroad, "Could not log in"),
		)

		// Hash the wrong password to prevent timing attacks
		_, _ = auth.SecurePassword(credentials.Password)

		return
	}

	// Get the user's password
	savedPassword, err := u.Tx.Queries().GetPassword(userId)
	if err != nil {
		u.Error(err.Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "Could not get password for user %v", userId.String()).
			Append(errors.LvlPlain, "Invalid credentials").
			Append(errors.LvlBroad, "Could not log in"),
		)

		// Hash the wrong password to prevent timing attacks
		_, _ = auth.SecurePassword(credentials.Password)

		return
	}

	// Verify the password
	if !auth.VerifyPassword(credentials.Password, savedPassword) {
		u.Error(errors.New().Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "Wrong password").
			Append(errors.LvlPlain, "Invalid credentials").
			Append(errors.LvlBroad, "Could not log in"),
		)
		return
	}

	// Silently update the user's password to a newer algorithm if applicable
	if !auth.PasswordStillSecure(savedPassword) {
		u.Logger.Infof("updating password %v for user to newer algorithm", credentials.Username)
		newPassword, err := auth.SecurePassword(credentials.Password)
		if err != nil {
			u.Error(err.
				Append(errors.LvlDebug, "Could not rehash password").
				Append(errors.LvlWordy, "Internal server error").
				Append(errors.LvlBroad, "Could not log in"),
			)
			return
		}
		err = u.Tx.Queries().UpdatePassword(userId, newPassword)
		if err != nil {
			u.Error(err.
				Append(errors.LvlDebug, "Could not update password").
				Append(errors.LvlWordy, "Database error").
				Append(errors.LvlBroad, "Could not log in"),
			)
			return
		}
	}

	// Create new session
	session := &types.Session{
		UserId:       userId,
		UserAgent:    c.Request.UserAgent(),
		IpAddress:    net.ParseIP(c.ClientIP()),
		IsShortLived: c.PostForm("remember") != "true",
		IsApi:        false,
	}
	err = u.Tx.Queries().InsertSession(session)
	if err != nil {
		u.Error(err.
			Append(errors.LvlBroad, "Could not log in"),
		)
		return
	}

	// Generate the token
	token, err := auth.NewToken(u.Config, u.Tx, userId, session.SessionId)
	if err != nil {
		u.Error(err.
			Append(errors.LvlWordy, "Could not generate token").
			Append(errors.LvlBroad, "Could not log in"),
		)
		return
	}

	u.Success(&gin.H{"token": token})
}

type registerPayload struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Email    string `form:"email"`
}

// TODO: check if registration is enabled on this instance otherwise we will
// TODO: have some kind of invite tokens that we will have to verify
func Register(c *gin.Context) {
	u := util.GetUtil(c)

	// Parse and validate the payload
	payload := registerPayload{}
	if err := c.ShouldBind(&payload); err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse payload").
			Append(errors.LvlWordy, "Malformed request").
			Append(errors.LvlBroad, "Could not register"),
		)
	}

	usernameErr := util.IsValidUsername(payload.Username)
	passwordErr := util.IsValidPassword(payload.Password)
	emailErr := util.IsValidEmail(payload.Email)
	if usernameErr != nil && passwordErr != nil && emailErr != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, usernameErr).AndErr(passwordErr).AndErr(emailErr).
			Append(errors.LvlDebug, "Input did not pass validation").
			Append(errors.LvlWordy, "Malformed request").
			Append(errors.LvlPlain, "Could not register"),
		)
		return
	}

	// Hash the password
	securedPassword, err := auth.SecurePassword(payload.Password)
	if err != nil {
		u.Error(err.
			Append(errors.LvlDebug, "Could not hash password").
			Append(errors.LvlWordy, "Internal server error").
			Append(errors.LvlBroad, "Could not register"),
		)
		return
	}

	// Check if any users exist to know if this user should be an admin
	usersExist, err := u.Tx.Queries().AnyUsersExist()
	if err != nil {
		u.Error(err.
			Append(errors.LvlDebug, "Could not check if any users exist").
			Append(errors.LvlWordy, "Database error").
			Append(errors.LvlBroad, "Could not register"),
		)
		return
	}

	// Construct the user
	user := &types.User{
		Username:       payload.Username,
		Email:          payload.Email,
		Admin:          !usersExist,
		Searchable:     true,
		ProfilePicture: util.GetGravatarUrl(payload.Email),
	}

	// Insert the user into the database
	userId, err := u.Tx.Queries().AddUser(user)
	if err != nil {
		u.Error(err.
			Append(errors.LvlBroad, "Could not register"),
		)
		return
	}

	// Initialize the user's settings
	err = u.Tx.Queries().InitializeUserSettings(userId)
	if err != nil {
		u.Error(err.
			Append(errors.LvlDebug, "Could not register"),
		)
		return
	}

	// Insert the password into the database
	err = u.Tx.Queries().InsertPassword(user.Id, securedPassword)
	if err != nil {
		u.Error(err.
			Append(errors.LvlDebug, "Could not insert password").
			Append(errors.LvlWordy, "Internal server error").
			Append(errors.LvlBroad, "Could not register"),
		)
		return
	}

	// Create new session
	session := &types.Session{
		UserId:       userId,
		UserAgent:    c.Request.UserAgent(),
		IpAddress:    net.ParseIP(c.ClientIP()),
		IsShortLived: c.PostForm("remember") != "true",
		IsApi:        false,
	}
	err = u.Tx.Queries().InsertSession(session)
	if err != nil {
		u.Error(err.
			Append(errors.LvlBroad, "Could not register"),
		)
		return
	}

	// Generate the token
	token, err := auth.NewToken(u.Config, u.Tx, userId, session.SessionId)
	if err != nil {
		u.Error(err.
			Append(errors.LvlWordy, "Could not generate token").
			Append(errors.LvlBroad, "Could not register"),
		)
	}

	u.Success(&gin.H{"token": token})
}

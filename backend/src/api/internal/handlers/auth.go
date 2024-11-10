package handlers

import (
	"fmt"
	"luna-backend/api/internal/context"
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// Parsing
	apiConfig := context.GetConfig(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	credentials := auth.BasicAuth{}
	if err := c.ShouldBind(&credentials); err != nil {
		apiConfig.Logger.Warn(err)
		util.Error(c, util.ErrorPayload)
		return
	}
	topErr := fmt.Errorf("failed to log in with credentials %v, %v", credentials.Username, credentials.Password)

	if util.IsValidUsername(credentials.Username) != nil || util.IsValidPassword(credentials.Password) != nil {
		apiConfig.Logger.Warnf("%v: user input failed validation", topErr)
		util.Error(c, util.ErrorPayload)
		return
	}

	// Check if the user exists
	userId, err := tx.Queries().GetUserIdFromUsername(credentials.Username)
	if err != nil {
		apiConfig.Logger.Warnf("%v: could not get user id for user %v: %v", topErr, credentials.Username, err)
		util.Error(c, util.ErrorInvalidCredentials)
		return
	}

	// Get the user's password
	savedPassword, err := tx.Queries().GetPassword(userId)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not get password for user %v: %v", topErr, credentials.Username, err)
		util.Error(c, util.ErrorInvalidCredentials)
		return
	}

	// Verify the password
	if !auth.VerifyPassword(credentials.Password, savedPassword) {
		apiConfig.Logger.Warnf("%v: passwords do not match", topErr)
		util.Error(c, util.ErrorInvalidCredentials)
		return
	}

	// Silently update the user's password to a newer algorithm if applicable
	if !auth.PasswordStillSecure(savedPassword) {
		apiConfig.Logger.Infof("updating password %v for user to newer algorithm", credentials.Username)
		newPassword, err := auth.SecurePassword(credentials.Password)
		if err != nil {
			apiConfig.Logger.Errorf("%v: could not hash password: %v", topErr, err)
			util.Error(c, util.ErrorInternal)
			return
		}
		err = tx.Queries().UpdatePassword(userId, newPassword)
		if err != nil {
			apiConfig.Logger.Errorf("%v: could not update password: %v", topErr, err)
			util.Error(c, util.ErrorDatabase)
			return
		}
	}

	// Generate the token
	token, err := auth.NewToken(apiConfig.CommonConfig, userId)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not generate token: %v", topErr, err)
		util.Error(c, util.ErrorInternal)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

type registerPayload struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Email    string `form:"email"`
}

// TODO: check if registration is enabled on this instance otherwise we will
// TODO: have some kind of invite tokens that we will have to verify
func Register(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	payload := registerPayload{}
	if err := c.ShouldBind(&payload); err != nil {
		apiConfig.Logger.Warn(err)
		util.Error(c, util.ErrorPayload)
		return
	}
	topErr := fmt.Errorf("failed to register user %v", payload.Username)

	if util.IsValidUsername(payload.Username) != nil || util.IsValidPassword(payload.Password) != nil || util.IsValidEmail(payload.Email) != nil {
		apiConfig.Logger.Warnf("%v: user input failed validation", topErr)
		util.Error(c, util.ErrorPayload)
		return
	}

	usersExist, err := tx.Queries().AnyUsersExist()
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not check if users exist: %v", topErr, err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	user := &types.User{
		Username: payload.Username,
		Email:    payload.Email,
		Admin:    !usersExist,
	}

	userId, err := tx.Queries().AddUser(user)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not add user: %v", topErr, err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	securedPassword, err := auth.SecurePassword(payload.Password)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not hash password: %v", topErr, err)
		util.Error(c, util.ErrorInternal)
		return
	}

	err = tx.Queries().InsertPassword(user.Id, securedPassword)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not insert password: %v", topErr, err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	// Generate the token
	token, err := auth.NewToken(apiConfig.CommonConfig, userId)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not generate token: %v", topErr, err)
		util.Error(c, util.ErrorInternal)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

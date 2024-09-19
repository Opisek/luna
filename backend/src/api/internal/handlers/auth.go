package handlers

import (
	"fmt"
	"luna-backend/api/internal/context"
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
		apiConfig.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "improper payload"})
		return
	}
	topErr := fmt.Errorf("failed to log in with credentials %v, %v", credentials.Username, credentials.Password)

	// Check if the user exists
	userId, err := tx.GetUserIdFromUsername(credentials.Username)
	if err != nil {
		apiConfig.Logger.Warnf("%v: could not get user id for user %v: %v", topErr, credentials.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Get the user's password
	savedPassword, algorithm, err := tx.GetPassword(userId)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not get password for user %v: %v", topErr, credentials.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify the password
	if !auth.VerifyPassword(credentials.Password, savedPassword, algorithm) {
		apiConfig.Logger.Warnf("%v: passwords do not match", topErr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Silently update the user's password to a newer algorithm if applicable
	if algorithm != auth.DefaultAlgorithm {
		apiConfig.Logger.Infof("updating password %v for user to newer algorithm", credentials.Username)
		hash, alg, err := auth.SecurePassword(credentials.Password)
		if err != nil {
			apiConfig.Logger.Errorf("%v: could not hash password: %v", topErr, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rehashing failed"})
			return
		}
		err = tx.UpdatePassword(userId, hash, alg)
		if err != nil {
			apiConfig.Logger.Errorf("%v: could not update password: %v", topErr, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rehashing failed"})
			return
		}
	}

	// Generate the token
	token, err := auth.NewToken(userId)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not generate token: %v", topErr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
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
		apiConfig.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "improper payload"})
		return
	}
	topErr := fmt.Errorf("failed to register user %v", payload.Username)

	hash, alg, err := auth.SecurePassword(payload.Password)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not hash password: %v", topErr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	usersExist, err := tx.AnyUsersExist()
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not check if users exist: %v", topErr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	user := &types.User{
		Username:  payload.Username,
		Password:  hash,
		Algorithm: alg,
		Email:     payload.Email,
		Admin:     !usersExist,
	}

	err = tx.AddUser(user)
	if err != nil {
		apiConfig.Logger.Errorf("%v: could not add user: %v", topErr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

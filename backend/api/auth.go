package api

import (
	"errors"
	"fmt"
	"luna-backend/auth"
	"luna-backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	credentials := auth.BasicAuth{}
	topErr := fmt.Errorf("failed to log in with credentials %v, %v", credentials.Username, credentials.Password)

	if err := c.ShouldBind(&credentials); err != nil {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "improper payload"})
		return
	}

	userId, err := apiConfig.db.GetUserIdFromUsername(credentials.Username)
	if err != nil {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	savedPassword, algorithm, err := apiConfig.db.GetPassword(userId)
	if err != nil {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !auth.VerifyPassword(credentials.Password, savedPassword, algorithm) {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.NewToken(credentials.Username)
	if err != nil {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

type registerPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// TODO: check if registration is enabled on this instance otherwise we will
// TODO: have some kind of invite tokens that we will have to verify
func register(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	payload := registerPayload{}
	topErr := fmt.Errorf("failed to register with payload %v", payload)

	if err := c.ShouldBind(&payload); err != nil {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "improper payload"})
		return
	}

	hash, alg, err := auth.SecurePassword(payload.Password)
	if err != nil {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	isFirstUser := !apiConfig.db.AnyUsersExist()

	user := &types.User{
		Username:  payload.Username,
		Password:  hash,
		Algorithm: alg,
		Email:     payload.Email,
		Admin:     isFirstUser,
	}

	err = apiConfig.db.AddUser(user)
	if err != nil {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")

		if err != nil || cookie == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token, err := auth.ParseToken(cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user", token.User)

		c.Next()
	}
}

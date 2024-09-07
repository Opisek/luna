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
	// Parsing
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	credentials := auth.BasicAuth{}
	if err := c.ShouldBind(&credentials); err != nil {
		apiConfig.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "improper payload"})
		return
	}
	topErr := fmt.Errorf("failed to log in with credentials %v, %v", credentials.Username, credentials.Password)

	// Check if the user exists
	userId, err := apiConfig.db.GetUserIdFromUsername(credentials.Username)
	if err != nil {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Get the user's password
	savedPassword, algorithm, err := apiConfig.db.GetPassword(userId)
	if err != nil {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify the password
	if !auth.VerifyPassword(credentials.Password, savedPassword, algorithm) {
		apiConfig.logger.Error(errors.Join(topErr, err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Silently update the user's password to a newer algorithm if applicable
	if algorithm != auth.DefaultAlgorithm {
		apiConfig.logger.Infof("updating password %v for user to newer algorithm", credentials.Username)
		hash, alg, err := auth.SecurePassword(credentials.Password)
		if err != nil {
			apiConfig.logger.Error(errors.Join(errors.New("could not update password"), err))
		}
		err = apiConfig.db.UpdatePassword(userId, hash, alg)
		if err != nil {
			apiConfig.logger.Error(errors.Join(errors.New("could not update password"), err))
		}
	}

	// Generate the token
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
	if err := c.ShouldBind(&payload); err != nil {
		apiConfig.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "improper payload"})
		return
	}
	topErr := fmt.Errorf("failed to register with payload %v", payload)

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

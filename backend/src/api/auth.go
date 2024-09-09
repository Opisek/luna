package api

import (
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
		apiConfig.logger.Errorf("%v: could not get user id for user %v: %v", topErr, credentials.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Get the user's password
	savedPassword, algorithm, err := apiConfig.db.GetPassword(userId)
	if err != nil {
		apiConfig.logger.Errorf("%v: could not get password for user %v: %v", topErr, credentials.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify the password
	if !auth.VerifyPassword(credentials.Password, savedPassword, algorithm) {
		apiConfig.logger.Errorf("%v: passwords do not match", topErr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Silently update the user's password to a newer algorithm if applicable
	if algorithm != auth.DefaultAlgorithm {
		apiConfig.logger.Infof("updating password %v for user to newer algorithm", credentials.Username)
		hash, alg, err := auth.SecurePassword(credentials.Password)
		if err != nil {
			apiConfig.logger.Errorf("%v: could not hash password: %v", topErr, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rehashing failed"})
			return
		}
		err = apiConfig.db.UpdatePassword(userId, hash, alg)
		if err != nil {
			apiConfig.logger.Errorf("%v: could not update password: %v", topErr, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rehashing failed"})
			return
		}
	}

	// Generate the token
	token, err := auth.NewToken(userId)
	if err != nil {
		apiConfig.logger.Errorf("%v: could not generate token: %v", topErr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
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
	topErr := fmt.Errorf("failed to register user %v", payload.Username)

	hash, alg, err := auth.SecurePassword(payload.Password)
	if err != nil {
		apiConfig.logger.Errorf("%v: could not hash password: %v", topErr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	usersExist, err := apiConfig.db.AnyUsersExist()
	if err != nil {
		apiConfig.logger.Errorf("%v: could not check if users exist: %v", topErr, err)
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

	err = apiConfig.db.AddUser(user)
	if err != nil {
		apiConfig.logger.Errorf("%v: could not add user: %v", topErr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieToken, cookieErr := c.Cookie("token")
		gotCookie := cookieErr == nil && cookieToken != ""
		bearerToken, bearerErr := getBearerToken(c)
		gotBearer := bearerErr == nil && bearerToken != ""

		if !gotCookie && !gotBearer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
			return
		}

		var token string
		if gotBearer {
			token = bearerToken
		} else if gotCookie {
			token = cookieToken
		}

		parsedToken, err := auth.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization token"})
			return
		}

		c.Set("user_id", parsedToken.UserId)

		c.Next()
	}
}

package api

import (
	"errors"
	"fmt"
	"luna-backend/auth"
	"luna-backend/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		apiConfig.logger.Error(errors.Join(
			topErr,
			errors.New("could not get user id"),
			err,
		))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Get the user's password
	savedPassword, algorithm, err := apiConfig.db.GetPassword(userId)
	if err != nil {
		apiConfig.logger.Error(errors.Join(
			topErr,
			errors.New("could not get password"),
			err,
		))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify the password
	if !auth.VerifyPassword(credentials.Password, savedPassword, algorithm) {
		apiConfig.logger.Error(errors.Join(
			topErr,
			errors.New("passwords do not match"),
			err,
		))
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
	token, err := auth.NewToken(userId)
	if err != nil {
		apiConfig.logger.Error(errors.Join(
			topErr,
			errors.New("could not generate token"),
			err,
		))
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
		apiConfig.logger.Error(errors.Join(
			topErr,
			errors.New("could not hash password"),
			err,
		))
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
		apiConfig.logger.Error(errors.Join(
			topErr,
			errors.New("could not add user"),
			err,
		))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func getBearerToken(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("missing bearer token")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return parts[1], nil
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

func getUserId(c *gin.Context) uuid.UUID {
	// it's fine to panic here because getUserId is always called after the
	// authMiddleware so we know the key must be set
	return c.MustGet("user_id").(uuid.UUID)
}

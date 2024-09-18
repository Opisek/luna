package handlers

import (
	"luna-backend/api/internal/context"
	"luna-backend/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieToken, cookieErr := c.Cookie("token")
		gotCookie := cookieErr == nil && cookieToken != ""
		bearerToken, bearerErr := context.GetBearerToken(c)
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

func TransactionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := context.GetConfig(c)
		tx, err := config.Db.BeginTransaction()

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		}

		c.Set("transaction", tx)

		c.Next()
	}
}

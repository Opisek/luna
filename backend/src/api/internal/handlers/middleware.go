package handlers

import (
	"luna-backend/api/internal/context"
	"luna-backend/api/internal/util"
	"luna-backend/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := context.GetConfig(c)

		cookieToken, cookieErr := c.Cookie("token")
		gotCookie := cookieErr == nil && cookieToken != ""
		bearerToken, bearerErr := context.GetBearerToken(c)
		gotBearer := bearerErr == nil && bearerToken != ""

		if !gotCookie && !gotBearer {
			util.Abort(c, util.ErrorTokenMissing)
			return
		}

		var token string
		if gotBearer {
			token = bearerToken
		} else if gotCookie {
			token = cookieToken
		}

		parsedToken, err := auth.ParseToken(config.CommonConfig, token)
		if err != nil {
			util.Abort(c, util.ErrorTokenInvalid)
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
			util.Abort(c, util.ErrorDatabase)
		}

		c.Set("transaction", tx)

		c.Next()
	}
}

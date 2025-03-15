package handlers

import (
	"context"
	ginContext "luna-backend/api/internal/context"
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"time"

	"github.com/gin-gonic/gin"
)

func ContextMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Set("context", ctx)

		finished := make(chan bool)
		go func() {
			c.Next()
			finished <- true
		}()

		select {
		case <-finished:
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				util.Abort(c, util.ErrorTimeout)
				ginContext.GetConfig(c).Logger.Errorf("request %v timed out after %v seconds", c.Request.URL.Path, timeout.Seconds())
			} else {
				util.Abort(c, util.ErrorUnknown)
				ginContext.GetConfig(c).Logger.Errorf("request %v failed: %v", c.Request.URL.Path, c.Err())
			}
		}
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := ginContext.GetConfig(c)

		cookieToken, cookieErr := c.Cookie("token")
		gotCookie := cookieErr == nil && cookieToken != ""
		bearerToken, bearerErr := ginContext.GetBearerToken(c)
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
		config := ginContext.GetConfig(c)
		tx, err := config.Db.BeginTransaction(ginContext.GetContext(c))

		if err != nil {
			config.Logger.Errorf("middleware failure: %v", err)
			util.Abort(c, util.ErrorDatabase)
		}

		c.Set("transaction", tx)

		c.Next()

		// If the handler did not explicitly commit the transaction, roll it back
		// just to be sure. Usually it should already be rolled back in the handler.
		//
		// A past bug in some handler caused the transaction to never be rolled
		// back, which caused the database to eventually stop responding. This
		// silently killed the server, because this middleware could not create any
		// nyw transcations (back then no timeout was being used either).
		//
		// Better to be safe than sorry and just roll back the transaction (again)
		// here, in the case another handler has an undiscovered bug.
		//
		// Edit: The context we now use automatically rolls back the transaction in
		// the context middleware. If the way we use contexts changes, we should to
		// reenable this rollback.
		//tx.Rollback(config.Logger)
	}
}

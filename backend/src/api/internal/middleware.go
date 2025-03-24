package middleware

import (
	"context"
	"fmt"
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/common"
	"luna-backend/db"
	"luna-backend/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestSetup(timeout time.Duration, database *db.Database, withTransaction bool, config *common.CommonConfig, logger *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseStatus := http.StatusOK
		var responseMsg *gin.H
		var responseFileName string
		var responseFileBody []byte
		var responseErr *errors.ErrorTrace
		var responseWarns []*errors.ErrorTrace

		// Final response sent at the end of the execution.
		defer func() {
			c.Header("Access-Control-Allow-Origin", config.Env.PUBLIC_URL)

			if responseFileBody != nil {
				c.Header("Content-Disposition", "attachment; filename="+responseFileName)
				c.Header("Content-Type", "application/text/plain")
				c.Header("Accept-Length", fmt.Sprintf("%d", len(responseFileBody)))
				var err error
				if c.Request.Method != http.MethodHead {
					_, err = c.Writer.Write(responseFileBody) // TODO: it would be nice if we could prevent the body from being fetched in first place
				}

				if err != nil {
					responseErr = errors.New().Status(http.StatusInternalServerError).
						AddErr(errors.LvlDebug, err).
						AltStr(errors.LvlPlain, "Could not download file")
				} else {
					return
				}
			}

			if responseErr != nil {
				logger.Error(responseErr.Serialize(errors.LvlDebug))
				c.AbortWithStatusJSON(responseErr.GetStatus(), &gin.H{"error": responseErr.Serialize(config.DetailLevel)})
				return
			}

			if responseMsg == nil {
				responseMsg = &gin.H{"status": "ok"}
			}

			if len(responseWarns) > 0 {
				warnStrs := make([]string, len(responseWarns))
				for i, warn := range responseWarns {
					logger.Warn(warn.Serialize(errors.LvlDebug))
					warnStrs[i] = warn.Serialize(config.DetailLevel)
				}

				(*responseMsg)["warnings"] = warnStrs
			}

			c.JSON(responseStatus, *responseMsg)
		}()

		// Timeout to be used by the handler and all its long-running functions (database queries, network request, ...)
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// If the request uses the database at all, we create a transaction for it.
		var tx *db.Transaction
		if withTransaction {
			tx, responseErr = database.BeginTransaction(ctx)
			if responseErr != nil {
				return
			}
		}

		// The handler will report back with a response body or an error trace
		responseChan := make(chan *util.Response)
		errChan := make(chan *errors.ErrorTrace)
		warnChan := make(chan *errors.ErrorTrace)

		// Pass important variables to the handler
		c.Set("transaction", tx)
		c.Set("config", config)
		c.Set("logger", logger)
		c.Set("context", ctx)

		// Pass important variables to the handler
		c.Set("transaction", tx)
		c.Set("config", config)
		c.Set("logger", logger)
		c.Set("context", ctx)

		c.Set("handlerUtil", &util.HandlerUtility{
			Config:       config,
			Logger:       logger,
			Tx:           tx,
			Context:      ctx,
			ResponseChan: responseChan,
			ErrChan:      errChan,
			WarnChan:     warnChan,
		})

		// Pass the execution on to the next middleware or the handler
		go func() {
			c.Next()
		}()

		// Gather warnings until we are done
		go func() {
			for {
				select {
				case warn := <-warnChan:
					responseWarns = append(responseWarns, warn)
				case <-ctx.Done():
					return
				}
			}
		}()

		// Wait until we either get a response, an error, or the time for the request expires
		select {
		// In case of a response
		case response := <-responseChan:
			responseStatus = response.GetStatus()
			responseMsg = response.GetMsg()
			responseFile := response.GetFile()

			if responseFile != nil {
				responseFileName = responseFile.GetName(tx.Queries())
				responseFileBody, responseErr = responseFile.GetBytes(tx.Queries())
				if responseErr != nil {
					return
				}
			}

			// Commit if the database was used
			if withTransaction {
				responseErr = tx.Commit(logger)
				if responseErr != nil {
					return
				}
			}

		// In case of a reported error
		case responseErr = <-errChan:
			// Rollback if the database was used
			if withTransaction {
				rollbackErr := tx.Rollback(logger)
				if rollbackErr != nil {
					logger.Error(rollbackErr.Serialize(errors.LvlDebug))
				}
			}

		// In case of a timeout or other unexpected error
		case <-ctx.Done():
			// Wait for a small time amount to see if we get a more detailed error about what exactly timed out
			select {
			case responseErr = <-errChan:
			case <-time.After(100 * time.Millisecond):
			}

			if responseErr == nil {
				responseErr = errors.New()
			}

			if ctx.Err() == context.DeadlineExceeded {
				// Would prefer to use StatusRequestTimeout, but then the browser keeps retrying indefinitely
				responseErr.Status(http.StatusGatewayTimeout).
					Append(errors.LvlWordy, "Request timed out after %v seconds", timeout.Seconds()).
					AltStr(errors.LvlBroad, "Request timed out")
			} else {
				responseErr.
					Append(errors.LvlBroad, "Request failed")
			}
		}
	}
}

func RequestAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := util.GetUtil(c)

		cookieToken, cookieErr := c.Cookie("token")
		gotCookie := cookieErr == nil && cookieToken != ""
		bearerToken, bearerErr := util.GetBearerToken(c)
		gotBearer := bearerErr == nil && bearerToken != ""

		if !gotCookie && !gotBearer {
			u.Error(errors.New().Status(http.StatusUnauthorized).
				Append(errors.LvlWordy, "Missing token"))
			c.Abort()
			return
		}

		var token string
		if gotBearer {
			token = bearerToken
		} else if gotCookie {
			token = cookieToken
		}

		parsedToken, err := auth.ParseToken(u.Config, token)
		if err != nil {
			u.Error(err)
			c.Abort()
			return
		}

		c.Set("user_id", parsedToken.UserId)

		c.Next()
	}
}

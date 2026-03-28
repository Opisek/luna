package middleware

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func DynamicThrottle(throttle *util.Throttle) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := util.GetUtil(c)

		ip := util.DetermineClientAddress(c).String()

		failCount := throttle.GetFailedAttempts(ip)
		if failCount > 50 {
			u.Logger.Warnf("IP %s has failed %d times in the last minute, aborting...", ip, failCount)
			u.Error(errors.New().Status(http.StatusTooManyRequests).
				Append(errors.LvlPlain, "Too many failed requests from this IP address"),
			)
			c.Abort()
		} else if failCount > 15 {
			u.Logger.Warnf("IP %s has failed %d times in the last minute, throttling...", ip, failCount)
			time.Sleep(5 * time.Second)
		} else if failCount > 10 {
			u.Logger.Warnf("IP %s has failed %d times in the last minute, throttling...", ip, failCount)
			time.Sleep(1 * time.Second)
		} else if failCount > 5 {
			u.Logger.Warnf("IP %s has failed %d times in the last minute, throttling...", ip, failCount)
			time.Sleep(100 * time.Millisecond)
		}

		c.Next()

		// Check for errors
		if _, errored := c.Get("error"); !errored {
			// Reset the number of failures if the request was successful
			throttle.RecordSuccessfulAttempt(ip)
		} else {
			// Increment the number of failures
			throttle.RecordFailedAttempt(ip)
		}
	}
}

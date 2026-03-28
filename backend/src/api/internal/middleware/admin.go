package middleware

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := util.GetUtil(c)
		userId := util.GetUserId(c)

		isAdmin, err := u.Tx.Queries().IsAdmin(userId)
		if err != nil {
			u.Error(err)
			c.Abort()
			return
		}

		if !isAdmin {
			u.Error(errors.New().Status(http.StatusForbidden).
				Append(errors.LvlPlain, "You must be an administrator to do this"))
			c.Abort()
			return
		}

		c.Next()
	}
}

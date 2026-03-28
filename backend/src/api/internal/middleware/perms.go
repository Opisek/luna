package middleware

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePermissions(requiredPerms ...types.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := util.GetUtil(c)
		permissions, exists := c.Get("permissions")
		if !exists {
			u.Error(errors.New().Status(http.StatusForbidden).
				Append(errors.LvlDebug, "No permissions found in context").
				AltStr(errors.LvlWordy, "You are missing one or more permissions").
				Append(errors.LvlPlain, "You are not authorized to perform this action"),
			)
			c.Abort()
			return
		}

		tokenPerms := permissions.(*types.TokenPermissions)

		for _, perm := range requiredPerms {
			if !tokenPerms.Has(perm) {
				u.Error(errors.New().Status(http.StatusForbidden).
					Append(errors.LvlDebug, "Missing permission: %v", perm).
					AltStr(errors.LvlWordy, "You are missing one or more permissions").
					Append(errors.LvlPlain, "You are not authorized to perform this action"),
				)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

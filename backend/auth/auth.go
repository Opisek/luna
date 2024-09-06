package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// WARNING
// This is not production code.
// This is only a quick setup to enable development.
// This code will be refactored for security in the future.
//

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")

		if err != nil || cookie == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token, err := ParseToken(cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user", token.User)

		c.Next()
	}
}

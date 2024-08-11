package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// WARNING
// This is not production code.
// This is only a quick setup to enable development.
// This code will be refactored for security in the future.
//

func Login(c *gin.Context) {
	credentials := BasicAuth{}

	if err := c.ShouldBind(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "improper payload"})
		return
	}

	fmt.Print(credentials.Password, credentials.Username)

	if credentials.Username != "admin" || credentials.Password != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := NewToken(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

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

package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": "admin",
	})

	signedToken, err := token.SignedString([]byte{'s', 'e', 'c', 'r', 'e', 't'})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}

package handlers

import (
	"luna-backend/api/internal/context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func GetVersion(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	c.JSON(http.StatusOK, gin.H{"version": apiConfig.CommonConfig.Version.String()})
}

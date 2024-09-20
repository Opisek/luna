package handlers

import (
	"luna-backend/api/internal/context"
	"luna-backend/api/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotImplemented(c *gin.Context) {
	util.Error(c, util.ErrorNotImplemented)
}

func GetVersion(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	c.JSON(http.StatusOK, gin.H{"version": apiConfig.CommonConfig.Version.String()})
}

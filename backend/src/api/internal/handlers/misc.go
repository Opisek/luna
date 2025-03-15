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

func GetHealth(c *gin.Context) {
	config := context.GetConfig(c)
	tx := context.GetTransaction(c)

	err := tx.Queries().CheckHealth()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		tx.Commit(config.Logger)
	} else {
		// With the current setup, this is never even reached, because the middleware already aborts the request earlier.
		// Stil, in the future we might have some other checks in CheckHealth.
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
		tx.Rollback(config.Logger)
	}
}

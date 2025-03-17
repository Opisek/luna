package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotImplemented(c *gin.Context) {
	u := util.GetUtil(c)
	u.Error(errors.New().Status(http.StatusNotImplemented))
}

func GetVersion(c *gin.Context) {
	u := util.GetUtil(c)
	u.Success(&gin.H{"version": u.Config.Version.String()})
}

func GetHealth(c *gin.Context) {
	u := util.GetUtil(c)

	err := u.Tx.Queries().CheckHealth()
	if err == nil {
		u.Success(&gin.H{"status": "ok"})
	} else {
		// With the current setup, this is never even reached, because the middleware already aborts the request earlier.
		// Stil, in the future we might have some other checks in CheckHealth.
		u.ResponseWithStatus(http.StatusInternalServerError, &gin.H{"status": "error"})
	}
}

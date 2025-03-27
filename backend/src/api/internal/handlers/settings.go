package handlers

import (
	"luna-backend/api/internal/util"

	"github.com/gin-gonic/gin"
)

func GetUserSettings(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)
	settings, err := u.Tx.Queries().GetRawUserSettings(userId)
	if err != nil {
		u.Error(err)
		return
	}

	u.SuccessRawJson(settings)
}

func GetGlobalSettings(c *gin.Context) {
	u := util.GetUtil(c)

	settings, err := u.Tx.Queries().GetRawGlobalSettings()
	if err != nil {
		u.Error(err)
		return
	}

	u.SuccessRawJson(settings)
}

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

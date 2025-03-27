package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"net/http"

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

func GetUserSetting(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)
	key := c.Param("settingKey")
	if key == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Missing setting name"))
		return
	}

	value, err := u.Tx.Queries().GetRawUserSetting(userId, key)
	if err != nil {
		u.Error(err)
		return
	}

	// TODO: would prefer { value: value } but we would need to unmarshal first
	u.SuccessRawJson(value)
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

func GetGlobalSetting(c *gin.Context) {
	u := util.GetUtil(c)

	key := c.Param("settingKey")
	if key == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Missing setting name"))
		return
	}

	value, err := u.Tx.Queries().GetRawGlobalSetting(key)
	if err != nil {
		u.Error(err)
		return
	}

	// TODO: would prefer { value: value } but we would need to unmarshal first
	u.SuccessRawJson(value)
}

package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/config"
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
			Append(errors.LvlWordy, "Missing setting key").
			AltStr(errors.LvlPlain, "Missing setting name"))
		return
	}

	setting, err := u.Tx.Queries().GetUserSetting(userId, key)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(&gin.H{"value": setting})
}

func PatchUserSetting(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)
	key := c.Param("settingKey")
	if key == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Missing setting key").
			AltStr(errors.LvlPlain, "Missing setting name"))
		return
	}

	setting, err := config.ParseUserSetting(key, []byte(c.PostForm("value")))
	if err != nil {
		u.Error(err)
		return
	}

	err = u.Tx.Queries().UpdateUserSetting(userId, setting)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

func ResetUserSetting(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)
	key := c.Param("settingKey")
	if key == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Missing setting key").
			AltStr(errors.LvlPlain, "Missing setting name"))
		return
	}

	setting, err := config.DefaultUserSetting(key)
	if err != nil {
		u.Error(err)
		return
	}

	err = u.Tx.Queries().UpdateUserSetting(userId, setting)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

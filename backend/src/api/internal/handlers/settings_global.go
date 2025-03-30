package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/config"
	"luna-backend/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
			Append(errors.LvlWordy, "Missing setting key").
			AltStr(errors.LvlPlain, "Missing setting name"))
		return
	}

	value, err := u.Tx.Queries().GetGlobalSetting(key)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(&gin.H{"value": value})
}

func PatchGlobalSetting(c *gin.Context) {
	u := util.GetUtil(c)

	key := c.Param("settingKey")
	if key == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Missing setting key").
			AltStr(errors.LvlPlain, "Missing setting name"))
		return
	}

	setting, err := config.ParseGlobalSetting(key, []byte(c.PostForm("value")))
	if err != nil {
		u.Error(err)
		return
	}

	u.Config.Settings.UpdateSetting(setting)

	err = u.Tx.Queries().UpdateGlobalSetting(setting)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

func ResetGlobalSetting(c *gin.Context) {
	u := util.GetUtil(c)

	key := c.Param("settingKey")
	if key == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Missing setting key").
			AltStr(errors.LvlPlain, "Missing setting name"))
		return
	}

	setting, err := config.DefaultGlobalSetting(key)
	if err != nil {
		u.Error(err)
		return
	}

	u.Config.Settings.UpdateSetting(setting)

	err = u.Tx.Queries().UpdateGlobalSetting(setting)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

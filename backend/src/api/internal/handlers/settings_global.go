package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/config"
	"luna-backend/constants"
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

func PatchGlobalSettings(c *gin.Context) {
	u := util.GetUtil(c)

	err := c.Request.ParseMultipartForm(constants.MaxFormBytes)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not parse form data").
			AltStr(errors.LvlPlain, "Malformed form data"))
		return
	}

	pairs := c.Request.PostForm

	if len(pairs) == 0 {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AltStr(errors.LvlPlain, "Nothing to change"))
		return
	}

	// We buffer all the settings into an array first, because it's faster to
	// parse than call the database. If one value is malformed, we save a lot of
	// time not running any queries.
	entries := make([]config.SettingsEntry, len(pairs))

	var tr *errors.ErrorTrace
	i := 0
	for key, value := range pairs {
		entries[i], tr = config.ParseGlobalSetting(key, []byte(value[0]))
		i++
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	// TODO: it may or may not be smarter to do INSERT ... ON CONFLICT than setting each key individually
	for _, setting := range entries {
		tr = u.Tx.Queries().UpdateGlobalSetting(setting)
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	for _, setting := range entries {
		u.Config.Settings.UpdateSetting(setting)
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

	err = u.Tx.Queries().UpdateGlobalSetting(setting)
	if err != nil {
		u.Error(err)
		return
	}

	u.Config.Settings.UpdateSetting(setting)

	u.Success(nil)
}

func ResetGlobalSettings(c *gin.Context) {
	u := util.GetUtil(c)

	settings := config.AllDefaultGlobalSettings()

	// TODO: it may or may not be smarter to do INSERT ... ON CONFLICT than setting each key individually
	for _, setting := range settings {
		err := u.Tx.Queries().UpdateGlobalSetting(setting)
		if err != nil {
			u.Error(err)
			return
		}
	}

	for _, setting := range settings {
		u.Config.Settings.UpdateSetting(setting)
	}

	u.Success(nil)
}

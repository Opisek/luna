package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/config"
	"luna-backend/constants"
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

func PatchUserSettings(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

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
			Append(errors.LvlPlain, "Nothing to change"))
		return
	}

	// We buffer all the settings into an array first, because it's faster to
	// parse than call the database. If one value is malformed, we save a lot of
	// time not running any queries.
	entries := make([]config.SettingsEntry, len(pairs))

	var tr *errors.ErrorTrace
	i := 0
	for key, value := range pairs {
		entries[i], tr = config.ParseUserSetting(key, []byte(value[0]))
		i++
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	// TODO: it may or may not be smarter to do INSERT ... ON CONFLICT than setting each key individually
	for _, setting := range entries {
		tr = u.Tx.Queries().UpdateUserSetting(userId, setting)
		if tr != nil {
			u.Error(tr)
			return
		}
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

func ResetUserSettings(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	settings := config.AllDefaultUserSettings()

	// TODO: it may or may not be smarter to do INSERT ... ON CONFLICT than setting each key individually
	for _, setting := range settings {
		err := u.Tx.Queries().UpdateUserSetting(userId, setting)
		if err != nil {
			u.Error(err)
			return
		}
	}

	u.Success(nil)
}

package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/config"
	"luna-backend/constants"
	"luna-backend/errors"
	"luna-backend/files"
	"luna-backend/types"
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
	//
	// Also, we take note of some settings that changed which might require
	// special treatment.
	entries := make([]config.SettingsEntry, len(pairs))

	gravatarPfpEnabled := u.Config.Settings.EnableGravatar.Enabled

	var tr *errors.ErrorTrace
	i := 0
	for key, value := range pairs {
		value, tr := config.ParseGlobalSetting(key, []byte(value[0]))
		entries[i] = value
		i++
		if tr != nil {
			u.Error(tr)
			return
		}

		if key == config.KeyEnableGravatar {
			gravatarPfpEnabled = value.(*config.EnableGravatar).Enabled
		}
	}

	for _, setting := range entries {
		// Update setting in the database
		// TODO: it may or may not be smarter to do INSERT ... ON CONFLICT than setting each key individually
		tr = u.Tx.Queries().UpdateGlobalSetting(setting)
		if tr != nil {
			u.Error(tr.
				Append(errors.LvlWordy, "Could not patch the setting %s in the database", setting.Key()).
				Append(errors.LvlPlain, "Could not update a global setting"),
			)
			return
		}

		// Additional actions after updating certain settings
		switch setting.Key() {

		case config.KeyEnableProfilePicturesUpload:
			if setting.(*config.EnableProfilePicturesUpload).Enabled {
				continue
			}
			// If profile picture uploads are disabled, delete all uploaded user profile pictures
			users, tx := u.Tx.Queries().GetUsers(true)
			if tx != nil {
				u.Error(tx.
					Append(errors.LvlWordy, "Could not perform action for setting %s", setting.Key()).
					Append(errors.LvlPlain, "Could not update a global setting"),
				)
				return
			}

			for _, user := range users {
				if user.ProfilePictureType != constants.ProfilePictureDatabase {
					continue
				}

				// Delete file
				dbFile := files.GetDatabaseFile(user.ProfilePictureFile)
				tx = u.Tx.Queries().DeleteFilecache(dbFile, user.Id)
				if tx != nil {
					u.Error(tx.
						Append(errors.LvlWordy, "Could not perform action for setting %s", setting.Key()).
						Append(errors.LvlPlain, "Could not update a global setting"),
					)
					return
				}

				// Reset profile picture
				user.ProfilePictureType = constants.ProfilePictureStatic
				user.ProfilePictureFile = types.EmptyId()
				user.ProfilePictureUrl = util.GetDefaultProfilePictureUrl(gravatarPfpEnabled, user.Email)
				tx = u.Tx.Queries().UpdateUserData(user)
				if tx != nil {
					u.Error(tx.
						Append(errors.LvlWordy, "Could not perform action for setting %s", setting.Key()).
						Append(errors.LvlPlain, "Could not update a global setting"),
					)
					return
				}
			}

		case config.KeyEnableGravatar:
			if setting.(*config.EnableGravatar).Enabled {
				continue
			}
			// If gravatar profile pictures are disabled, delete all gravatar profile pictures
			users, tx := u.Tx.Queries().GetUsers(true)
			if tx != nil {
				u.Error(tx.
					Append(errors.LvlWordy, "Could not perform action for setting %s", setting.Key()).
					Append(errors.LvlPlain, "Could not update a global setting"),
				)
				return
			}

			for _, user := range users {
				if user.ProfilePictureType != constants.ProfilePictureGravatar {
					continue
				}

				// Delete file if exists
				if !user.ProfilePictureFile.IsEmpty() {
					dbFile := files.GetDatabaseFile(user.ProfilePictureFile)
					tx = u.Tx.Queries().DeleteFilecache(dbFile, user.Id)
					if tx != nil {
						u.Error(tx.
							Append(errors.LvlWordy, "Could not perform action for setting %s", setting.Key()).
							Append(errors.LvlPlain, "Could not update a global setting"),
						)
						return
					}
				}

				// Reset profile picture
				user.ProfilePictureType = constants.ProfilePictureStatic
				user.ProfilePictureFile = types.EmptyId()
				user.ProfilePictureUrl = util.GetDefaultProfilePictureUrl(false, user.Email)
				tx = u.Tx.Queries().UpdateUserData(user)
				if tx != nil {
					u.Error(tx.
						Append(errors.LvlWordy, "Could not perform action for setting %s", setting.Key()).
						Append(errors.LvlPlain, "Could not update a global setting"),
					)
					return
				}
			}
		default:
		}
	}

	// Update settings in the config
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

	// TODO: can we do this in bulk rather than individually for each key?
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

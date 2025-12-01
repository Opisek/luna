package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/constants"
	"luna-backend/errors"
	"luna-backend/files"
	"luna-backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)
	user, err := u.Tx.Queries().GetUser(userId)
	if err != nil {
		u.Error(err)
		return
	}
	err = user.UpdateEffectiveProfilePicture(u.Config.Settings.CacheProfilePictures.Enabled)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(&gin.H{"user": user})
}

func GetUsers(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	all := c.Query("all") == "true"
	if all && !util.HasAdminPrivilegesAndReportError(c) {
		return
	}

	users, tr := u.Tx.Queries().GetUsers(all)
	if tr != nil {
		u.Error(tr)
		return
	}

	for _, user := range users {
		tr = user.UpdateEffectiveProfilePicture(u.Config.Settings.CacheProfilePictures.Enabled)
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	if all {
		u.Success(&gin.H{
			"users":   users,
			"current": userId,
		})
	} else {
		// When all is not set, i.e., the request comes from a non-admin user,
		// cast User to StrippedUser to remove (potentially) sensitive information.
		strippedUsers := make([]*types.StrippedUser, 0, len(users))
		for _, user := range users {
			strippedUser := &types.StrippedUser{
				Id:                         user.Id,
				Username:                   user.Username,
				Admin:                      user.Admin,
				EffectiveProfilePictureUrl: user.EffectiveProfilePictureUrl,
			}
			strippedUsers = append(strippedUsers, strippedUser)
		}
		u.Success(&gin.H{
			"users":   strippedUsers,
			"current": userId,
		})
	}
}

func PatchUserData(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	// Parse data to change
	newUsername := c.PostForm("username")
	newEmail := c.PostForm("email")
	newPassword := c.PostForm("new_password")
	newProfilePictureType := c.PostForm("pfp_type")
	rawNewProfilePictureUrl := c.PostForm("pfp_url")
	newPfpFileHeader, pfpFileErr := c.FormFile("pfp_file")
	newProfilePictureRequested := newProfilePictureType != "" || rawNewProfilePictureUrl != "" || pfpFileErr == nil
	rawNewSearchable := c.PostForm("searchable")

	switch pfpFileErr {
	case nil:
		break
	case http.ErrMissingFile:
		newPfpFileHeader = nil
	default:
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, pfpFileErr).
			Append(errors.LvlPlain, "Invalid form data"))
	}

	if newUsername == "" && newEmail == "" && newPassword == "" && !newProfilePictureRequested {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Nothing to change"))
		return
	}

	if newUsername != "" && util.IsValidUsername(newUsername) != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid username"))
		return
	}
	if newEmail != "" && util.IsValidEmail(newEmail) != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid email"))
		return
	}
	if newPassword != "" && util.IsValidPassword(newPassword) != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid password"))
		return
	}

	newSearchable := false
	if rawNewSearchable != "" {
		if rawNewSearchable != "true" && rawNewSearchable != "false" {
			u.Error(errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Invalid searchable value"))
			return
		}
		newSearchable = rawNewSearchable == "true"
	}

	oldUserStruct, tr := u.Tx.Queries().GetUser(userId)
	if tr != nil {
		u.Error(tr)
		return
	}

	// Set new profile picture
	var newProfilePictureUrl *types.Url
	newProfilePictureId := types.EmptyId()
	if newProfilePictureRequested {
		// Delete old profile picture if applicable
		oldPfpFile := *files.GetDatabaseFile(oldUserStruct.ProfilePictureFile)
		if oldPfpFile.GetId() != types.EmptyId() {
			tr = u.Tx.Queries().DeleteFilecache(&oldPfpFile, userId)
			if tr != nil {
				u.Error(tr.
					Append(errors.LvlDebug, "Could not delete old profile picture file %v", oldPfpFile.GetId()).
					Append(errors.LvlPlain, "Could not update profile picture"))
				return
			}
		}

		// Refresh profile picture if the email changes and gravatar is used
		if !newProfilePictureRequested && newEmail != "" && oldUserStruct.ProfilePictureType == "gravatar" {
			newProfilePictureType = "gravatar"
			currentGravatarUrl, err := types.NewUrl(oldUserStruct.ProfilePictureUrl.String())
			if err != nil {
				u.Error(errors.New().Status(http.StatusInternalServerError).
					AddErr(errors.LvlDebug, err).
					Append(errors.LvlWordy, "Could not parse old gravatar profile picture").
					Append(errors.LvlPlain, "Could not update profile picture"),
				)
				return
			}
			newProfilePictureUrl = util.GetGravatarUrlWithParams(newEmail, currentGravatarUrl.URL().RawQuery)
		}

		// Parse new profile picture
		switch newProfilePictureType {
		case constants.ProfilePictureGravatar:
			fallthrough
		case constants.ProfilePictureStatic:
			fallthrough
		case constants.ProfilePictureRemote:
			if rawNewProfilePictureUrl == "" {
				u.Error(errors.New().Status(http.StatusBadRequest).
					Append(errors.LvlPlain, "Missing profile picture URL"))
				return
			}

			var err error
			err = util.IsValidUrl(rawNewProfilePictureUrl)
			if err != nil {
				u.Error(errors.New().Status(http.StatusBadRequest).
					Append(errors.LvlPlain, "Invalid profile picture URL").
					AddErr(errors.LvlWordy, err))
				return
			}
			newProfilePictureUrl, err = types.NewUrl(rawNewProfilePictureUrl)
			if err != nil {
				u.Error(errors.New().Status(http.StatusBadRequest).
					Append(errors.LvlPlain, "Invalid profile picture URL").
					AddErr(errors.LvlWordy, err))
				return
			}

		case constants.ProfilePictureDatabase:
			if newPfpFileHeader == nil {
				u.Error(errors.New().Status(http.StatusBadRequest).
					Append(errors.LvlPlain, "Missing profile picture file"))
				return
			}

			pfpFile, err := newPfpFileHeader.Open()
			if err != nil {
				u.Error(errors.New().Status(http.StatusBadRequest).
					AddErr(errors.LvlDebug, err).
					Append(errors.LvlPlain, "Could not open profile picture file"))
				return
			}

			uploadedFile, tr := files.NewDatabaseFileFromContent(newPfpFileHeader.Filename, pfpFile, userId, u.Tx.Queries())
			if tr != nil {
				u.Error(tr.
					Append(errors.LvlDebug, "Could not create file from content").
					Append(errors.LvlPlain, "Could not upload profile picture"))
				return
			}

			newProfilePictureId = uploadedFile.GetId()

		default:
			u.Error(errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Invalid profile picture type"))
			return
		}

		// Create local profile picture cache
		if newProfilePictureType == constants.ProfilePictureGravatar || newProfilePictureType == constants.ProfilePictureRemote {
			pfpFile, tr := files.NewRemoteFile(newProfilePictureUrl, "image/*", auth.NewNoAuth(), userId, u.Tx.Queries())
			if tr != nil {
				u.Error(tr.
					Append(errors.LvlDebug, "Could not create remote file for profile picture").
					Append(errors.LvlPlain, "Could not upload profile picture"),
				)
				return
			}

			newProfilePictureId = pfpFile.GetId()
		}
	}
	// Reauthenticate if needed
	reauthenticationRequired := newUsername != "" || newEmail != "" || newPassword != ""

	if reauthenticationRequired {
		password := c.PostForm("password")
		if password == "" {
			u.Error(errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Missing password"))
			return
		}

		// Get the user's password
		savedPassword, err := u.Tx.Queries().GetPassword(userId)
		if err != nil {
			u.Error(err.Status(http.StatusUnauthorized).
				Append(errors.LvlDebug, "Could not get password for user %v", userId.String()).
				Append(errors.LvlPlain, "Invalid credentials"),
			)
			return
		}

		// Verify the password
		if !auth.VerifyPassword(password, savedPassword, u.Config) {
			u.Error(errors.New().Status(http.StatusUnauthorized).
				Append(errors.LvlDebug, "Wrong password").
				Append(errors.LvlPlain, "Invalid credentials"),
			)
			return
		}
	}

	// Update the user
	var newUserStruct *types.User
	if newUsername != "" || newEmail != "" || rawNewSearchable != "" || newProfilePictureRequested {
		newUserStruct = &types.User{
			Id:                 userId,
			Username:           newUsername,
			Email:              newEmail,
			Searchable:         newSearchable,
			ProfilePictureType: newProfilePictureType,
			ProfilePictureUrl:  newProfilePictureUrl,
			ProfilePictureFile: newProfilePictureId,
		}

		newUserStruct.Admin = oldUserStruct.Admin

		if newUsername == "" {
			newUserStruct.Username = oldUserStruct.Username
		}
		if newEmail == "" {
			newUserStruct.Email = oldUserStruct.Email
		}
		if rawNewSearchable == "" {
			newUserStruct.Searchable = oldUserStruct.Searchable
		}
		if !newProfilePictureRequested {
			newUserStruct.ProfilePictureType = oldUserStruct.ProfilePictureType
			newUserStruct.ProfilePictureUrl = oldUserStruct.ProfilePictureUrl
			newUserStruct.ProfilePictureFile = oldUserStruct.ProfilePictureFile
		}

		tr = u.Tx.Queries().UpdateUserData(newUserStruct)
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	// Update the password
	if newPassword != "" {
		securedPassword, err := auth.SecurePassword(newPassword, u.Config)
		if err != nil {
			u.Error(err.
				Append(errors.LvlDebug, "Could not hash new password"),
			)
			return
		}
		err = u.Tx.Queries().UpdatePassword(userId, securedPassword)
		if err != nil {
			u.Error(err.
				Append(errors.LvlDebug, "Could not update password"),
			)
			return
		}
	}

	response := &gin.H{
		"status": "ok",
	}

	if newProfilePictureType != "" {
		tr = newUserStruct.UpdateEffectiveProfilePicture(u.Config.Settings.CacheProfilePictures.Enabled)
		if tr != nil {
			u.Error(tr)
			return
		}
		(*response)["profile_picture"] = newUserStruct.EffectiveProfilePictureUrl.String()
	}

	u.Success(response)
}

func DeleteUser(c *gin.Context) {
	u := util.GetUtil(c)

	executingUserId := util.GetUserId(c)
	affectedUserId, tr := util.GetIdOrDefault(c, "user", "self", executingUserId)
	if tr != nil {
		u.Error(tr)
		return
	}
	if affectedUserId != executingUserId && !util.HasAdminPrivilegesAndReportError(c) {
		return
	}

	// Get the user's password
	password := c.PostForm("password")
	if password == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing password"),
		)
		return
	}

	// Get the user's password
	savedPassword, err := u.Tx.Queries().GetPassword(executingUserId)
	if err != nil {
		u.Error(err.Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "Could not get password for user %v", executingUserId.String()).
			Append(errors.LvlPlain, "Invalid credentials"),
		)
		return
	}

	// Verify the password
	if !auth.VerifyPassword(password, savedPassword, u.Config) {
		u.Error(errors.New().Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "Wrong password").
			Append(errors.LvlPlain, "Invalid credentials"),
		)
		return
	}

	err = u.Tx.Queries().DeleteUser(affectedUserId)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

func EnableUser(c *gin.Context) {
	u := util.GetUtil(c)

	userId, tr := util.GetId(c, "user")
	if tr != nil {
		u.Error(tr)
		return
	}

	tr = u.Tx.Queries().SetUserEnabled(userId, true)

	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(nil)
}

func DisableUser(c *gin.Context) {
	u := util.GetUtil(c)

	userId, tr := util.GetId(c, "user")
	if tr != nil {
		u.Error(tr)
		return
	}

	isAdmin, tr := u.Tx.Queries().IsAdmin(userId)
	if tr != nil {
		u.Error(tr)
		return
	}
	if isAdmin {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Admin accounts cannot be disabled."),
		)
		return
	}

	tr = u.Tx.Queries().SetUserEnabled(userId, false)

	if tr != nil {
		u.Error(tr)
		return
	}

	tr = u.Tx.Queries().DeleteSessions(userId)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(nil)
}

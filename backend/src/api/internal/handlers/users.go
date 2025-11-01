package handlers

import (
	"fmt"
	"luna-backend/api/internal/util"
	"luna-backend/auth"
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

	u.Success(&gin.H{"user": user})
}

func GetUsers(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	all := c.Query("all") == "true"
	if all {
		isAdmin, tr := u.Tx.Queries().IsAdmin(userId)
		if tr != nil {
			u.Error(tr)
			return
		}
		if !isAdmin {
			u.Error(errors.New().Status(http.StatusForbidden).
				Append(errors.LvlPlain, "Only administrators can view all users"),
			)
			return
		}
	}

	users, tr := u.Tx.Queries().GetUsers(all)
	if tr != nil {
		u.Error(tr)
		return
	}

	// TODO: when not using all=true, we might want to hide some fields like verified, email, etc.

	u.Success(&gin.H{
		"users":   users,
		"current": userId,
	})
}

func PatchUserData(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	// Parse data to change
	newUsername := c.PostForm("username")
	newEmail := c.PostForm("email")
	newPassword := c.PostForm("new_password")
	rawNewProfilePicture := c.PostForm("pfp_url")
	pfpFileHeader, pfpFileErr := c.FormFile("pfp_file")
	rawNewSearchable := c.PostForm("searchable")

	switch pfpFileErr {
	case nil:
		break
	case http.ErrMissingFile:
		pfpFileHeader = nil
	default:
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, pfpFileErr).
			Append(errors.LvlPlain, "Invalid form data"))
	}

	if newUsername == "" && newEmail == "" && newPassword == "" && rawNewProfilePicture == "" && rawNewSearchable == "" && pfpFileHeader == nil {
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

	var newProfilePicture *types.Url
	if rawNewProfilePicture != "" {
		var err error
		err = util.IsValidUrl(rawNewProfilePicture)
		if err != nil {
			u.Error(errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Invalid profile picture URL").
				AddErr(errors.LvlWordy, err))
			return
		}
		newProfilePicture, err = types.NewUrl(rawNewProfilePicture)
		if err != nil {
			u.Error(errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Invalid profile picture URL").
				AddErr(errors.LvlWordy, err))
			return
		}
	} else if pfpFileHeader != nil {
		pfpFile, err := pfpFileHeader.Open()
		if err != nil {
			u.Error(errors.New().Status(http.StatusBadRequest).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlPlain, "Could not open profile picture file"))
			return
		}

		uploadedFile, tr := files.NewDatabaseFileFromContent(pfpFileHeader.Filename, pfpFile, userId, u.Tx.Queries())
		if tr != nil {
			u.Error(tr.
				Append(errors.LvlDebug, "Could not create file from content").
				Append(errors.LvlPlain, "Could not upload profile picture"))
			return
		}

		fileId := uploadedFile.GetId()
		fileUrl := fmt.Sprintf("/api/files/%s", fileId.String())
		newProfilePicture, err = types.NewUrl(fileUrl)
		if err != nil {
			u.Error(errors.New().Status(http.StatusInternalServerError).
				Append(errors.LvlDebug, "Could not create file url").
				Append(errors.LvlPlain, "Could not upload profile picture"))
			return
		}
	}

	// Delete old profile picture if applicable
	oldUserStruct, tr := u.Tx.Queries().GetUser(userId)
	if tr != nil {
		u.Error(tr)
		return
	}

	if newProfilePicture != nil {
		// Check if the old profile picture is a database file
		oldFileId, err := util.IsDatabaseFileUrl(oldUserStruct.ProfilePicture)
		// Delete if it is
		if err == nil {
			oldFile := files.GetDatabaseFile(oldFileId)
			tr := u.Tx.Queries().DeleteFilecache(oldFile, userId)
			if tr != nil {
				tr = tr.Append(errors.LvlWordy, "Could not delete old profile picture")
				u.Warn(tr)
			}
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
		if !auth.VerifyPassword(password, savedPassword) {
			u.Error(errors.New().Status(http.StatusUnauthorized).
				Append(errors.LvlDebug, "Wrong password").
				Append(errors.LvlPlain, "Invalid credentials"),
			)
			return
		}
	}

	// Update the user
	if newUsername != "" || newEmail != "" || rawNewSearchable != "" || newProfilePicture != nil {
		newUserStruct := &types.User{
			Id:             userId,
			Username:       newUsername,
			Email:          newEmail,
			Searchable:     newSearchable,
			ProfilePicture: newProfilePicture,
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
		if newProfilePicture == nil {
			newUserStruct.ProfilePicture = oldUserStruct.ProfilePicture
		}

		tr = u.Tx.Queries().UpdateUserData(newUserStruct)
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	// Update the password
	if newPassword != "" {
		securedPassword, err := auth.SecurePassword(newPassword)
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

	if newProfilePicture != nil {
		(*response)["profile_picture"] = newProfilePicture.String()
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
	if affectedUserId != executingUserId {
		isAdmin, tr := u.Tx.Queries().IsAdmin(executingUserId)
		if tr != nil {
			u.Error(tr)
			return
		}
		if !isAdmin {
			u.Error(errors.New().Status(http.StatusForbidden).
				Append(errors.LvlPlain, "You are not allowed to delete other user accounts"),
			)
			return
		}
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
	if !auth.VerifyPassword(password, savedPassword) {
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

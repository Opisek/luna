package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserData(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)
	user, err := u.Tx.Queries().GetUserData(userId)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(&gin.H{"user": user})
}

func PatchUserData(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	// Parse the request
	newUsername := c.PostForm("username")
	newEmail := c.PostForm("email")
	newPassword := c.PostForm("new_password")
	rawNewProfilePicture := c.PostForm("profile_picture")
	rawNewSearchable := c.PostForm("searchable")

	if newUsername == "" && newEmail == "" && newPassword == "" && rawNewProfilePicture == "" && rawNewSearchable == "" {
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
		if util.IsValidUrl(rawNewProfilePicture) != nil {
			u.Error(errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Invalid profile picture URL"))
			return
		}
		var err error
		newProfilePicture, err = types.NewUrl(rawNewProfilePicture)
		if err != nil {
			u.Error(errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Invalid profile picture url"))
			return
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
	if newUsername != "" || newEmail != "" || rawNewSearchable != "" || rawNewProfilePicture != "" {
		newUserStruct := &types.User{
			Id:             userId,
			Username:       newUsername,
			Email:          newEmail,
			Searchable:     newSearchable,
			ProfilePicture: newProfilePicture,
		}

		oldUserStruct, tr := u.Tx.Queries().GetUserData(userId)
		if tr != nil {
			u.Error(tr)
			return
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
		if rawNewProfilePicture == "" {
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

	u.Success(nil)
}

func DeleteUser(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	// Get the user's password
	password := c.PostForm("password")
	if password == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing password"),
		)
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

	err = u.Tx.Queries().DeleteUser(userId)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

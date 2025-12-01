package types

import (
	"fmt"
	"luna-backend/constants"
	"luna-backend/errors"
	"net/http"
	"time"
)

type User struct {
	Id ID `json:"id"`

	Username string `json:"username"`
	Email    string `json:"email"`

	Admin      bool `json:"admin"`
	Verified   bool `json:"verified"`
	Enabled    bool `json:"enabled"`
	Searchable bool `json:"searchable"`

	ProfilePictureType         string `json:"profile_picture_type"`
	ProfilePictureUrl          *Url   `json:"profile_picture_url"`
	ProfilePictureFile         ID     `json:"profile_picture_file"`
	EffectiveProfilePictureUrl *Url   `json:"profile_picture"`

	CreatedAt time.Time `json:"created_at"`
}

type StrippedUser struct {
	Id                         ID     `json:"id"`
	Username                   string `json:"username"`
	Admin                      bool   `json:"admin"`
	EffectiveProfilePictureUrl *Url   `json:"profile_picture"`
}

// Profile picture types:
// - Remote: Set url to the remote url and file to the cached file id.
// - Database: Set url to nil and file to the database file id.
// - Gravatar: Set url to gravatar url and file to nil.
// - Static: Set url to the static (relative) file url and file to nil.

func (user *User) UpdateEffectiveProfilePicture(cacheProfirePicures bool) *errors.ErrorTrace {
	fromDatabase := user.ProfilePictureType == constants.ProfilePictureDatabase || (cacheProfirePicures && (user.ProfilePictureType == constants.ProfilePictureGravatar || user.ProfilePictureType == constants.ProfilePictureRemote))

	if fromDatabase {
		fileUrl := fmt.Sprintf("/api/files/%s", user.ProfilePictureFile.String())
		newProfilePictureUrl, err := NewUrl(fileUrl)
		if err != nil {
			return errors.New().Status(http.StatusInternalServerError).
				Append(errors.LvlDebug, "Could not create file url for user %s with file ID %s", user.Id, user.ProfilePictureFile).
				Append(errors.LvlPlain, "Could not get profile picture")
		}
		user.EffectiveProfilePictureUrl = newProfilePictureUrl
	} else {
		user.EffectiveProfilePictureUrl = user.ProfilePictureUrl
	}

	return nil
}

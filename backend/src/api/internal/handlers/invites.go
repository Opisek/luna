package handlers

import (
	"fmt"
	"luna-backend/api/internal/util"
	"luna-backend/constants"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/files"
	"luna-backend/types"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func GetInvites(c *gin.Context) {
	u := util.GetUtil(c)

	// Get the invites
	invites, tr := u.Tx.Queries().GetValidInvites()
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{
		"invites": invites,
	})
}

func GetInviteQrCode(c *gin.Context) {
	u := util.GetUtil(c)

	// Invite ID
	inviteId, tr := util.GetId(c, "invite")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Get the invite code
	invite, tr := u.Tx.Queries().GetValidInviteById(inviteId)
	if tr != nil {
		u.Error(tr)
		return
	}

	// Generate the invite link
	inviteLink := fmt.Sprintf("%s/register?code=%s", u.Config.Env.PUBLIC_URL, invite.Code)

	// Generate the QR code
	qrCode, err := qrcode.Encode(inviteLink, qrcode.Medium, 256)
	if err != nil {
		u.Error(errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Could not generate QR code"),
		)
		return
	}

	u.ResponseWithFile(files.NewVolatileFile(fmt.Sprintf("invite-%s.png", invite.Code), qrCode))
}

func PutInvite(c *gin.Context) {
	u := util.GetUtil(c)

	// Invite author
	userId := util.GetUserId(c)

	// Invitee email
	// Optional, can either be a valid email address or empty
	email := c.PostForm("email")
	if email != "" {
		if err := util.IsValidEmail(email); err != nil {
			u.Error(errors.New().Status(http.StatusBadRequest).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlPlain, "Invalid email address"),
			)
			return
		}
	}

	// Calculate duration
	durationString := c.PostForm("duration")
	if durationString == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing duration"),
		)
		return
	}

	duration, err := strconv.Atoi(durationString)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Invalid duration"),
		)
		return
	}

	if duration < 0 {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlDebug, "Duration cannot be negative").
			Append(errors.LvlPlain, "Invalid duration"),
		)
		return
	}

	if duration > int(constants.MaxInviteDuration.Seconds()) {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlDebug, "Duration cannot be greater than %d hours", (int)(constants.MaxInviteDuration.Hours())).
			Append(errors.LvlPlain, "Invalid duration"),
		)
		return
	}

	currentTime := time.Now()

	var tr *errors.ErrorTrace

	// Invite code
	// 36 possible characters (0-9, A-Z)
	// 3 groups of 4 characters
	// 36^16 = approx 62 bits of entropy
	// This is far enough considering the request throttling that is also in place.
	// If we want this to be more secure, 4 groups would result in approx 83 bits of entropy.
	var code string
	for code == "" || strings.Contains(code, "+") || strings.Contains(code, "/") {
		random, tr := crypto.GenerateRandomBase64(16)
		if tr != nil {
			u.Error(tr.Status(http.StatusInternalServerError).
				Append(errors.LvlWordy, "Could not generate random invite code"),
			)
		}
		code = strings.ToUpper(fmt.Sprintf("%s-%s-%s", random[:4], random[4:8], random[8:12]))
	}

	// Create invite
	invite := &types.RegistrationInvite{
		Author:    userId,
		Email:     email,
		CreatedAt: currentTime,
		Expires:   currentTime.Add(time.Duration(duration) * time.Second),
		Code:      code,
	}

	tr = u.Tx.Queries().InsertInvite(invite)
	if tr != nil {
		u.Error(tr.Status(http.StatusInternalServerError).
			Append(errors.LvlWordy, "Could not insert invite"),
		)
		return
	}

	// TODO: Send email if address provided

	u.Success(&gin.H{
		"invite": invite,
	})
}

func DeleteInvite(c *gin.Context) {
	u := util.GetUtil(c)

	// Invite ID
	invite, tr := util.GetId(c, "invite")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Delete invite
	tr = u.Tx.Queries().DeleteInvite(invite)
	if tr != nil {
		u.Error(tr.Status(http.StatusInternalServerError).
			Append(errors.LvlWordy, "Could not delete invite"),
		)
		return
	}

	u.Success(nil)
}

func DeleteInvites(c *gin.Context) {
	u := util.GetUtil(c)

	// Delete all invites
	tr := u.Tx.Queries().DeleteInvites()
	if tr != nil {
		u.Error(tr.Status(http.StatusInternalServerError).
			Append(errors.LvlWordy, "Could not delete invites"),
		)
		return
	}

	u.Success(nil)
}

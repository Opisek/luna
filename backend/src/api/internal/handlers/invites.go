package handlers

import (
	"fmt"
	"luna-backend/api/internal/util"
	"luna-backend/constants"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func PutInvite(c *gin.Context) {
	u := util.GetUtil(c)

	// Invite author
	userId := util.GetUserId(c)

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

	// Invite code
	random, tr := crypto.GenerateRandomBase64(16)
	if tr != nil {
		u.Error(tr.Status(http.StatusInternalServerError).
			Append(errors.LvlWordy, "Could not generate random invite code"),
		)
	}

	code := strings.ToUpper(fmt.Sprintf("%s-%s-%s-%s", random[:4], random[4:8], random[8:12], random[12:16]))

	// Create invite
	invite := &types.RegistrationInvite{
		Author:  userId,
		Expires: currentTime.Add(time.Duration(duration) * time.Second),
		Code:    code,
	}

	tr = u.Tx.Queries().InsertInvite(invite)
	if tr != nil {
		u.Error(tr.Status(http.StatusInternalServerError).
			Append(errors.LvlWordy, "Could not insert invite"),
		)
		return
	}

	u.Success(&gin.H{
		"code": code,
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

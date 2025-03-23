package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"luna-backend/files"
	"luna-backend/types"
	"net/http"
	"strings"

	"github.com/emersion/go-webdav/caldav"
	"github.com/gin-gonic/gin"
)

func NotImplemented(c *gin.Context) {
	u := util.GetUtil(c)
	u.Error(errors.New().Status(http.StatusNotImplemented))
}

func GetVersion(c *gin.Context) {
	u := util.GetUtil(c)
	u.Success(&gin.H{"version": u.Config.Version.String()})
}

func GetHealth(c *gin.Context) {
	u := util.GetUtil(c)

	err := u.Tx.Queries().CheckHealth()
	if err == nil {
		u.Success(&gin.H{"status": "ok"})
	} else {
		// With the current setup, this is never even reached, because the middleware already aborts the request earlier.
		// Stil, in the future we might have some other checks in CheckHealth.
		u.ResponseWithStatus(http.StatusInternalServerError, &gin.H{"status": "error"})
	}
}

func CheckUrl(c *gin.Context) {
	u := util.GetUtil(c)

	// Pares URL aund auth parameters
	auth, tr := parseAuthMethod(c)
	if tr != nil {
		u.Error(tr)
		return
	}
	rawUrl := c.PostForm("url")
	if rawUrl == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing url"))
		return
	}
	if util.IsValidUrl(rawUrl) != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid url"))
		return
	}
	url, err := types.NewUrl(rawUrl)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Invalid url"))
	}

	isIcal, tr, statusCode := isUrlIcal(u, url, auth)
	if tr != nil {
		u.Error(tr)
		return
	}
	if isIcal {
		u.Success(&gin.H{
			"type": types.SourceIcal,
		})
		return
	}
	if statusCode == http.StatusUnauthorized {
		u.Success(&gin.H{
			"type":   types.SourceUnknown,
			"status": statusCode,
		})
		return
	}

	isCaldav, tr, principalUrl := isUrlCaldav(u, url, auth)
	if tr != nil {
		if strings.Contains(tr.Serialize(errors.LvlDebug), "401 Unauthorized") {
			statusCode = http.StatusUnauthorized
		}
		if statusCode != http.StatusOK {
			u.Success(&gin.H{
				"type":   types.SourceUnknown,
				"status": statusCode,
			})
		} else {
			u.Error(tr)
		}
		return
	}
	if isCaldav {
		u.Success(&gin.H{
			"type": types.SourceCaldav,
			"url":  principalUrl,
		})
		return
	}

	u.Success(&gin.H{
		"type":   types.SourceUnknown,
		"status": statusCode,
	})
}

func isUrlIcal(u *util.HandlerUtility, url *types.Url, auth types.AuthMethod) (bool, *errors.ErrorTrace, int) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not create request").
			Append(errors.LvlWordy, "Could not check url %v", url).
			AltStr(errors.LvlPlain, "Could not check url"), 0
	}

	req = req.WithContext(u.Context)

	res, err := auth.Do(req)
	if err != nil {
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not check url %v", url), 0
	}

	if res.StatusCode != http.StatusOK {
		return false, nil, res.StatusCode
	}

	tr := files.IsValidIcalFile(res.Body, u.Tx.Queries())
	if tr != nil {
		return false, tr, http.StatusOK
	}

	return true, nil, http.StatusOK
}

func isUrlCaldav(u *util.HandlerUtility, url *types.Url, auth types.AuthMethod) (bool, *errors.ErrorTrace, string) {
	client, err := caldav.NewClient(
		auth,
		url.String(),
	)
	if err != nil {
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not create CalDAV client"), ""
	}

	principalUrl, err := client.FindCurrentUserPrincipal(u.Context)
	if err != nil {
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get current user principal"), ""
	}

	if principalUrl == "" || strings.HasPrefix(principalUrl, "/") {
		principalUrl = url.URL().Scheme + "://" + url.URL().Host + principalUrl
	}

	return true, nil, principalUrl
}

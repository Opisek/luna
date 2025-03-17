package handlers

import (
	"net/http"

	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/errors"
	"luna-backend/interface/primitives"
	"luna-backend/interface/protocols/caldav"
	"luna-backend/interface/protocols/ical"
	"luna-backend/types"

	"github.com/gin-gonic/gin"
)

type exposedSource struct {
	Id   types.ID `json:"id"`
	Name string   `json:"name"`
	Type string   `json:"type"`
}

type exposedDetailedSource struct {
	Id       types.ID    `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Settings interface{} `json:"settings"`
	AuthType string      `json:"auth_type"`
	Auth     interface{} `json:"auth"`
}

func getSources(u *util.HandlerUtility, userId types.ID) ([]primitives.Source, *errors.ErrorTrace) {
	srcs, err := u.Tx.Queries().GetSourcesByUser(userId)
	if err != nil {
		return nil, err
	}
	return srcs, nil
}

func GetSources(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	sources, err := getSources(u, userId)
	if err != nil {
		u.Error(err)
		return
	}

	exposedSources := make([]exposedSource, len(sources))
	for i, source := range sources {
		exposedSources[i] = exposedSource{
			Id:   source.GetId(),
			Name: source.GetName(),
			Type: source.GetType(),
		}
	}

	u.Success(&gin.H{"sources": exposedSources})
}

func GetSource(c *gin.Context) {
	u := util.GetUtil(c)

	sourceId, err := util.GetId(c, "source")
	if err != nil {
		u.Error(err)
		return
	}

	userId := util.GetUserId(c)

	source, err := u.Tx.Queries().GetSource(userId, sourceId)
	if err != nil {
		u.Error(err)
		return
	}

	exposedSource := exposedDetailedSource{
		Id:       source.GetId(),
		Name:     source.GetName(),
		Settings: source.GetSettings(),
		Type:     source.GetType(),
		AuthType: source.GetAuth().GetType(),
		Auth:     source.GetAuth(),
	}

	u.Success(&gin.H{"source": exposedSource})
}

func parseAuthMethod(c *gin.Context) (auth.AuthMethod, *errors.ErrorTrace) {
	var sourceAuth auth.AuthMethod

	authType := c.PostForm("auth_type")
	switch authType {
	case types.AuthNone:
		sourceAuth = auth.NewNoAuth()
	case types.AuthBasic:
		username := c.PostForm("auth_username")
		password := c.PostForm("auth_password")
		if username == "" || password == "" {
			return nil, errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Missing username or password")
		}

		sourceAuth = auth.NewBasicAuth(username, password)
	case types.AuthBearer:
		token := c.PostForm("auth_token")
		if token == "" {
			return nil, errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Missing token")
		}

		sourceAuth = auth.NewBearerAuth(token)
	case "":
		return nil, errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing authentication type")
	default:
		return nil, errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Unknown authentication type: %v", authType)
	}

	return sourceAuth, nil
}

func parseSource(c *gin.Context, sourceName string, sourceAuth auth.AuthMethod, q types.DatabaseQueries) (primitives.Source, *errors.ErrorTrace) {
	var source primitives.Source

	sourceType := c.PostForm("type")
	switch sourceType {
	case types.SourceCaldav:
		rawUrl := c.PostForm("url")
		if rawUrl == "" {
			return nil, errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Missing CalDAV url")
		}
		if util.IsValidUrl(rawUrl) != nil {
			return nil, errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Invalid CalDAV url")
		}
		sourceUrl, err := types.NewUrl(rawUrl)
		if err != nil {
			return nil, errors.New().Status(http.StatusBadRequest).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlPlain, "Invalid CalDAV url")
		}

		source = caldav.NewCaldavSource(sourceName, sourceUrl, sourceAuth)
	case types.SourceIcal:
		locationType := c.PostForm("location")
		if locationType == "" {
			return nil, errors.New().Status(http.StatusBadRequest).
				Append(errors.LvlPlain, "Missing iCal location")
		}

		switch locationType {
		case "remote":
			rawUrl := c.PostForm("url")
			if rawUrl == "" {
				return nil, errors.New().Status(http.StatusBadRequest).
					Append(errors.LvlPlain, "Missing iCal url")
			}
			if util.IsValidUrl(rawUrl) != nil {
				return nil, errors.New().Status(http.StatusBadRequest).
					Append(errors.LvlPlain, "Invalid iCal url")
			}
			sourceUrl, err := types.NewUrl(rawUrl)
			if err != nil {
				return nil, errors.New().Status(http.StatusBadRequest).
					AddErr(errors.LvlDebug, err).
					Append(errors.LvlPlain, "Invalid iCal url")
			}
			source = ical.NewRemoteIcalSource(sourceName, sourceUrl, sourceAuth)
		case "local":
			if sourceAuth.GetType() != types.AuthNone {
				return nil, errors.New().Status(http.StatusBadRequest).
					Append(errors.LvlPlain, "Local iCal sources do not support authentication")
			}
			rawPath := c.PostForm("path")
			if rawPath == "" {
				return nil, errors.New().Status(http.StatusBadRequest).
					Append(errors.LvlPlain, "Missing iCal path")
			}
			sourcePath, err := types.NewPath(rawPath)
			if err != nil {
				return nil, errors.New().Status(http.StatusBadRequest).
					AddErr(errors.LvlDebug, err).
					Append(errors.LvlPlain, "Invalid iCal path")
			}
			source = ical.NewLocalIcalSource(sourceName, sourcePath)
		case "database":
			if sourceAuth.GetType() != types.AuthNone {
				return nil, errors.New().Status(http.StatusBadRequest).
					Append(errors.LvlPlain, "Database iCal sources do not support authentication")
			}

			fileHeader, err := c.FormFile("file")
			if err != nil || fileHeader == nil {
				return nil, errors.New().Status(http.StatusBadRequest).
					AddErr(errors.LvlDebug, err).
					Append(errors.LvlPlain, "Missing or corrupted iCal file")
			}
			if fileHeader.Size > 50*1000*1000 {
				return nil, errors.New().Status(http.StatusInsufficientStorage).
					Append(errors.LvlPlain, "iCal file too large")
			}

			file, err := fileHeader.Open()
			if err != nil {
				return nil, errors.New().Status(http.StatusInternalServerError).
					AddErr(errors.LvlDebug, err).
					Append(errors.LvlPlain, "Could not open iCal file")
			}

			var tr *errors.ErrorTrace
			source, tr = ical.NewDatabaseIcalSource(sourceName, fileHeader.Filename, file, q)
			if tr != nil {
				return nil, tr
			}
		}

	case "":
		return nil, errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing source type")
	default:
		return nil, errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Unknown source type: %v", sourceType)
	}

	return source, nil
}

func PutSource(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	sourceName := c.PostForm("name")
	if sourceName == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Missing name"))
		return
	}

	sourceAuth, err := parseAuthMethod(c)
	if err != nil {
		u.Error(err)
		return
	}

	source, err := parseSource(c, sourceName, sourceAuth, u.Tx.Queries())
	if err != nil {
		u.Error(err)
		return
	}

	id, err := u.Tx.Queries().InsertSource(userId, source)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(&gin.H{"id": id.String()})
}

func PatchSource(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	sourceId, err := util.GetId(c, "source")
	if err != nil {
		u.Error(err)
		return
	}

	newName := c.PostForm("name")
	newType := c.PostForm("type")
	newAuthType := c.PostForm("auth_type")

	if newName == "" && newType == "" && newAuthType == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Nothing to change"))
		return
	}

	source, err := u.Tx.Queries().GetSource(userId, sourceId)
	if err != nil {
		u.Error(err)
		return
	}

	var newAuth auth.AuthMethod = nil
	if newAuthType != "" {
		newAuth, err = parseAuthMethod(c)
		if err != nil {
			u.Error(err)
			return
		}
	}

	var newSourceSettings primitives.SourceSettings = nil
	if newType != "" {
		if newAuth == nil {
			newAuth = source.GetAuth()
		}
		newSource, err := parseSource(c, newName, newAuth, u.Tx.Queries())
		if err != nil {
			u.Error(err)
			return
		}
		newSourceSettings = newSource.GetSettings()
	}
	if source.GetType() == "ical" {
		if newType == "ical" {
			err = source.Cleanup(u.Tx.Queries())
		}
		if err != nil {
			u.Error(err.
				Append(errors.LvlWordy, "Could not clean up source before editing"))
			return
		}
	}

	err = u.Tx.Queries().UpdateSource(userId, sourceId, newName, newAuth, newType, newSourceSettings)
	if err != nil {
		u.Error(err)
		return
	}

	u.Success(nil)
}

func DeleteSource(c *gin.Context) {
	u := util.GetUtil(c)

	userId := util.GetUserId(c)

	sourceId, err := util.GetId(c, "source")
	if err != nil {
		u.Error(err)
		return
	}

	source, err := u.Tx.Queries().GetSource(userId, sourceId)
	if err != nil {
		u.Warn(err)
	} else {
		err = source.Cleanup(u.Tx.Queries())
		if err != nil {
			u.Warn(err.
				Append(errors.LvlWordy, "Could not clean up source before deleting"))
		}
	}

	deleted, err := u.Tx.Queries().DeleteSource(userId, sourceId)
	if err != nil {
		u.Error(err)
		return
	}

	if deleted {
		u.Success(nil)
	} else {
		u.Error(errors.New().Status(http.StatusNotFound).
			Append(errors.LvlPlain, "Source not found"))
	}
}

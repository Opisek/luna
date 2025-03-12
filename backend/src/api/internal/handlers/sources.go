package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"luna-backend/api/internal/config"
	"luna-backend/api/internal/context"
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/db"
	"luna-backend/interface/primitives"
	"luna-backend/interface/protocols/caldav"
	"luna-backend/interface/protocols/ical"
	"luna-backend/types"

	"github.com/gin-gonic/gin"
)

type exposedSource struct {
	Id   types.ID `json:"id"`
	Name string   `json:"name"`
}

type exposedDetailedSource struct {
	Id       types.ID    `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Settings interface{} `json:"settings"`
	AuthType string      `json:"auth_type"`
	Auth     interface{} `json:"auth"`
}

func getSources(_ *config.Api, tx *db.Transaction, userId types.ID) ([]primitives.Source, error) {
	srcs, err := tx.Queries().GetSources(userId)
	if err != nil {
		return nil, fmt.Errorf("could not get sources: %v", err)
	}
	return srcs, nil
}

func GetSources(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	sources, err := getSources(apiConfig, tx, userId)
	if err != nil {
		apiConfig.Logger.Error(err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	exposedSources := make([]exposedSource, len(sources))
	for i, source := range sources {
		exposedSources[i] = exposedSource{
			Id:   source.GetId(),
			Name: source.GetName(),
		}
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusOK, exposedSources)
}

func GetSource(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	sourceId, err := context.GetId(c, "source")
	if err != nil {
		apiConfig.Logger.Errorf("could not get source id: %v", err)
		util.Error(c, util.ErrorMalformedID)
		return
	}

	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	source, err := tx.Queries().GetSource(userId, sourceId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get source: %v", err)
		util.Error(c, util.ErrorDatabase)
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

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusOK, exposedSource)
}

func parseAuthMethod(c *gin.Context) (auth.AuthMethod, error) {
	var sourceAuth auth.AuthMethod

	authType := c.PostForm("auth_type")
	switch authType {
	case types.AuthNone:
		sourceAuth = auth.NewNoAuth()
	case types.AuthBasic:
		username := c.PostForm("auth_username")
		password := c.PostForm("auth_password")
		if username == "" || password == "" {
			return nil, errors.New("missing username or password")
		}

		sourceAuth = auth.NewBasicAuth(username, password)
	case types.AuthBearer:
		token := c.PostForm("auth_token")
		if token == "" {
			return nil, errors.New("missing token")
		}

		sourceAuth = auth.NewBearerAuth(token)
	case "":
		return nil, errors.New("missing auth type")
	default:
		return nil, errors.New("invalid auth type")
	}

	return sourceAuth, nil
}

func parseSource(c *gin.Context, sourceName string, sourceAuth auth.AuthMethod, q types.DatabaseQueries) (primitives.Source, error) {
	var source primitives.Source

	sourceType := c.PostForm("type")
	switch sourceType {
	case types.SourceCaldav:
		rawUrl := c.PostForm("url")
		if rawUrl == "" {
			return nil, errors.New("missing caldav url")
		}
		if util.IsValidUrl(rawUrl) != nil {
			return nil, errors.New("invalid caldav url")
		}
		sourceUrl, err := types.NewUrl(rawUrl)
		if err != nil {
			return nil, errors.New("invalid caldav url")
		}

		source = caldav.NewCaldavSource(sourceName, sourceUrl, sourceAuth)
	case types.SourceIcal:
		locationType := c.PostForm("location")
		if locationType == "" {
			return nil, errors.New("missing ical location")
		}

		switch locationType {
		case "remote":
			rawUrl := c.PostForm("url")
			if rawUrl == "" {
				return nil, errors.New("missing ical url")
			}
			if util.IsValidUrl(rawUrl) != nil {
				return nil, errors.New("invalid ical url")
			}
			sourceUrl, err := types.NewUrl(rawUrl)
			if err != nil {
				return nil, errors.New("invalid ical url")
			}
			source = ical.NewRemoteIcalSource(sourceName, sourceUrl, sourceAuth)
		case "local":
			if sourceAuth.GetType() != types.AuthNone {
				return nil, errors.New("local ical sources cannot have auth")
			}
			rawPath := c.PostForm("path")
			if rawPath == "" {
				return nil, errors.New("missing ical path")
			}
			sourcePath, err := types.NewPath(rawPath)
			if err != nil {
				return nil, errors.New("invalid ical path")
			}
			source = ical.NewLocalIcalSource(sourceName, sourcePath)
		case "database":
			if sourceAuth.GetType() != types.AuthNone {
				return nil, errors.New("database ical sources cannot have auth")
			}

			fileHeader, err := c.FormFile("file")
			if err != nil || fileHeader == nil {
				return nil, fmt.Errorf("missing or errornous ical file: %w", err)
			}
			if fileHeader.Size > 50*1000*1000 {
				return nil, errors.New("ical file too large")
			}

			file, err := fileHeader.Open()
			if err != nil {
				return nil, fmt.Errorf("could not open ical file: %w", err)
			}

			source, err = ical.NewDatabaseIcalSource(sourceName, file, q)
			if err != nil {
				return nil, fmt.Errorf("could not create database ical source: %w", err)
			}
		}

	case "":
		return nil, errors.New("invalid source type")
	default:
		return nil, errors.New("missing source type")
	}

	return source, nil
}

func PutSource(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	sourceName := c.PostForm("name")
	if sourceName == "" {
		apiConfig.Logger.Warn("missing name")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailName)
		return
	}

	sourceAuth, err := parseAuthMethod(c)
	if err != nil {
		apiConfig.Logger.Warnf("could not parse auth: %v", err)
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailAuth)
		return
	}

	source, err := parseSource(c, sourceName, sourceAuth, tx.Queries())
	if err != nil {
		apiConfig.Logger.Warnf("could not parse source: %v", err)
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailSource)
		return
	}

	id, err := tx.Queries().InsertSource(userId, source)
	if err != nil {
		apiConfig.Logger.Errorf("could not insert source %v for user %v: %v", source.GetId().String(), userId.String(), err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id.String()})
}

func PatchSource(c *gin.Context) {
	var err error

	apiConfig := context.GetConfig(c)
	sourceId, err := context.GetId(c, "source")
	if err != nil {
		apiConfig.Logger.Errorf("could not get source id: %v", err)
		util.Error(c, util.ErrorMalformedID)
		return
	}

	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	newName := c.PostForm("name")
	newType := c.PostForm("type")
	newAuthType := c.PostForm("auth_type")

	if newName == "" && newType == "" && newAuthType == "" {
		apiConfig.Logger.Warn("no values to change")
		util.ErrorDetailed(c, util.ErrorPayload, util.DetailFields)
		return
	}

	var newAuth auth.AuthMethod = nil
	if newAuthType != "" {
		newAuth, err = parseAuthMethod(c)
		if err != nil {
			apiConfig.Logger.Warnf("could not parse auth: %v", err)
			util.ErrorDetailed(c, util.ErrorPayload, util.DetailAuth)
			return
		}
	}

	var newSourceSettings primitives.SourceSettings = nil
	if newType != "" {
		newSource, err := parseSource(c, newName, newAuth, tx.Queries())
		if err != nil {
			apiConfig.Logger.Warnf("could not parse source: %v", err)
			util.ErrorDetailed(c, util.ErrorPayload, util.DetailSource)
			return
		}
		newSourceSettings = newSource.GetSettings()
	}

	source, err := tx.Queries().GetSource(userId, sourceId)
	if err != nil {
		apiConfig.Logger.Errorf("could not get source: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	} else if source.GetType() == "ical" {
		err = source.Cleanup(tx.Queries())
		if err != nil {
			apiConfig.Logger.Errorf("error cleaning up source before editing: %v", err)
			util.Error(c, util.ErrorDatabase)
			return
		}
	}

	err = tx.Queries().UpdateSource(userId, sourceId, newName, newAuth, newType, newSourceSettings)
	if err != nil {
		apiConfig.Logger.Errorf("could not update source: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	util.Success(c)
}

func DeleteSource(c *gin.Context) {
	apiConfig := context.GetConfig(c)
	sourceId, err := context.GetId(c, "source")
	if err != nil {
		apiConfig.Logger.Errorf("could not get source id: %v", err)
		util.Error(c, util.ErrorMalformedID)
		return
	}

	userId := context.GetUserId(c)
	tx := context.GetTransaction(c)
	defer tx.Rollback(apiConfig.Logger)

	source, err := tx.Queries().GetSource(userId, sourceId)
	if err != nil {
		apiConfig.Logger.Warnf("could not get source: %v", err)
	} else {
		err = source.Cleanup(tx.Queries())
		if err != nil {
			apiConfig.Logger.Warnf("error cleaning up after source: %v", err)
		}
	}

	deleted, err := tx.Queries().DeleteSource(userId, sourceId)
	if err != nil {
		apiConfig.Logger.Errorf("could not delete source: %v", err)
		util.Error(c, util.ErrorDatabase)
		return
	}

	if tx.Commit(apiConfig.Logger) != nil {
		util.Error(c, util.ErrorDatabase)
		return
	}

	if deleted {
		util.Success(c)
	} else {
		util.Error(c, util.ErrorSourceNotFound)
	}
}

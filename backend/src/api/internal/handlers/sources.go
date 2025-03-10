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

func parseSource(c *gin.Context, sourceName string, sourceAuth auth.AuthMethod) (primitives.Source, error) {
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
		fallthrough
	default:
		return nil, errors.New("invalid source type")
	case "":
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

	source, err := parseSource(c, sourceName, sourceAuth)
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
		newSource, err := parseSource(c, newName, newAuth)
		if err != nil {
			apiConfig.Logger.Warnf("could not parse source: %v", err)
			util.ErrorDetailed(c, util.ErrorPayload, util.DetailSource)
			return
		}
		newSourceSettings = newSource.GetSettings()
	}

	apiConfig.Logger.Debugf("parsed params")

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

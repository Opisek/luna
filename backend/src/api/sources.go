package api

import (
	"errors"
	"net/http"

	"luna-backend/auth"
	"luna-backend/sources"
	"luna-backend/sources/caldav"
	"luna-backend/types"

	"github.com/gin-gonic/gin"
)

type exposedSource struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type exposedDetailedSource struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Settings interface{} `json:"settings"`
	AuthType string      `json:"auth_type"`
	Auth     interface{} `json:"auth"`
}

func getSources(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	userId := getUserId(c)

	sources, err := apiConfig.db.GetSources(userId)
	if err != nil {
		apiConfig.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get sources"})
		return
	}

	exposedSources := make([]exposedSource, len(sources))
	for i, source := range sources {
		exposedSources[i] = exposedSource{
			Id:   source.GetId().String(),
			Name: source.GetName(),
		}
	}

	c.JSON(http.StatusOK, exposedSources)
}

func getSource(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	userId := getUserId(c)
	sourceId, err := getSourceId(c)
	if err != nil {
		apiConfig.logger.Errorf("could not get source id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed or missing source id"})
		return
	}

	source, err := apiConfig.db.GetSource(userId, sourceId)
	if err != nil {
		apiConfig.logger.Errorf("could not get source: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get source"})
		return
	}

	exposedSource := exposedDetailedSource{
		Id:       source.GetId().String(),
		Name:     source.GetName(),
		Settings: source.GetSettings(),
		Type:     source.GetType(),
		AuthType: source.GetAuth().GetType(),
		Auth:     source.GetAuth(),
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

func parseSource(c *gin.Context, sourceName string, sourceAuth auth.AuthMethod) (sources.Source, error) {
	var source sources.Source

	sourceType := c.PostForm("type")
	switch sourceType {
	case types.SourceCaldav:
		rawUrl := c.PostForm("url")
		if rawUrl == "" {
			return nil, errors.New("missing caldav url")
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

func putSource(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	userId := getUserId(c)

	sourceName := c.PostForm("name")
	if sourceName == "" {
		apiConfig.logger.Error("missing name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing name"})
		return
	}

	sourceAuth, err := parseAuthMethod(c)
	if err != nil {
		apiConfig.logger.Errorf("could not parse auth: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source, err := parseSource(c, sourceName, sourceAuth)
	if err != nil {
		apiConfig.logger.Errorf("could not parse source: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := apiConfig.db.InsertSource(userId, source)
	if err != nil {
		apiConfig.logger.Errorf("could not insert source %v for user %v: %v", source.GetId().String(), userId.String(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not insert source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id.String()})
}

func patchSource(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	var err error

	userId := getUserId(c)
	sourceId, err := getSourceId(c)
	if err != nil {
		apiConfig.logger.Errorf("could not get source id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed or missing source id"})
		return
	}

	newName := c.PostForm("name")
	newType := c.PostForm("type")
	newAuthType := c.PostForm("auth_type")

	if newName == "" && newType == "" && newAuthType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	var newAuth auth.AuthMethod = nil
	if newAuthType != "" {
		newAuth, err = parseAuthMethod(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "malformed auth"})
			return
		}
	}

	var newSourceSettings sources.SourceSettings = nil
	if newType != "" {
		newSource, err := parseSource(c, newName, newAuth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newSourceSettings = newSource.GetSettings()
	}

	apiConfig.logger.Debugf("parsed params")

	err = apiConfig.db.UpdateSource(userId, sourceId, newName, newAuth, newType, newSourceSettings)
	if err != nil {
		apiConfig.logger.Errorf("could not update source: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func deleteSource(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	userId := getUserId(c)
	sourceId, err := getSourceId(c)
	if err != nil {
		apiConfig.logger.Errorf("could not get source id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed or missing source id"})
		return
	}

	err = apiConfig.db.DeleteSource(userId, sourceId)
	if err != nil {
		apiConfig.logger.Errorf("could not delete source: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

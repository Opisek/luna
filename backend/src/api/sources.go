package api

import (
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

func putSource(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	userId := getUserId(c)

	var source sources.Source

	sourceName := c.PostForm("name")
	if sourceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing name"})
		return
	}

	var sourceAuth auth.AuthMethod

	authType := c.PostForm("auth_type")
	switch authType {
	case auth.AuthNone:
		sourceAuth = auth.NewNoAuth()
	case auth.AuthBasic:
		username := c.PostForm("auth_username")
		password := c.PostForm("auth_password")
		if username == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing username or password"})
			return
		}

		sourceAuth = auth.NewBasicAuth(username, password)
	case auth.AuthBearer:
		token := c.PostForm("auth_token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
			return
		}

		sourceAuth = auth.NewBearerAuth(token)
	case "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing auth type"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auth type"})
		return
	}

	sourceType := c.PostForm("type")
	switch sourceType {
	case sources.SourceCaldav:
		rawUrl := c.PostForm("url")
		if rawUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing caldav url"})
			return
		}
		sourceUrl, err := types.NewUrl(rawUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid caldav url"})
			return
		}

		source = caldav.NewCaldavSource(sourceName, sourceUrl, sourceAuth)
	case sources.SourceIcal:
		fallthrough
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source type"})
		return
	case "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing source type"})
		return
	}

	id, err := apiConfig.db.InsertSource(userId, source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id.String()})
}

func patchSource(c *gin.Context) {
	notImplemented(c)

	//apiConfig := getConfig(c)
	//if apiConfig == nil {
	//	return
	//}

	//userId := getUserId(c)
	//rawSourceId := c.Param("sourceId")

	//if rawSourceId == "" {
	//	apiConfig.logger.Error("could not patch source: missing source id")
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "missing source id"})
	//	return
	//}
	//sourceId, err := sources.SourceIdFromString(rawSourceId)
	//if err != nil {
	//	apiConfig.logger.Errorf("could not patch source: malformed source id: %v", err)
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "malformed source id"})
	//	return
	//}

	//var source sources.Source
}

func deleteSource(c *gin.Context) {
	apiConfig := getConfig(c)
	if apiConfig == nil {
		return
	}

	userId := getUserId(c)
	rawSourceId := c.Param("sourceId")

	if rawSourceId == "" {
		apiConfig.logger.Error("could not delete source: missing source id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing source id"})
		return
	}
	sourceId, err := types.IdFromString(rawSourceId)
	if err != nil {
		apiConfig.logger.Errorf("could not delete source: malformed source id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed source id"})
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

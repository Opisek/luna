package api

import (
	"net/http"
	"net/url"

	"luna-backend/auth"
	"luna-backend/sources"
	"luna-backend/sources/caldav"

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

	c.JSON(http.StatusOK, gin.H{"sources": sources})
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
		sourceUrl, err := url.Parse(rawUrl)
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

	err := apiConfig.db.InsertSource(userId, source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add source"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": source.GetId().String()})
}

func patchSource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func deleteSource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

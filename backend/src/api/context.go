package api

import (
	"errors"
	"fmt"
	"luna-backend/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getConfig(c *gin.Context) *Api {
	// TODO: consider changing to "MustGet"
	apiConfig, err := c.Get("apiConfig")
	if !err {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "context error"})
		return nil
	}
	return apiConfig.(*Api)
}

func getUserId(c *gin.Context) types.ID {
	// it's fine to panic here because getUserId is always called after the
	// authMiddleware so we know the key must be set
	return c.MustGet("user_id").(types.ID)
}

func getSourceId(c *gin.Context) (types.ID, error) {
	rawSourceId := c.Param("sourceId")

	if rawSourceId == "" {
		return types.EmptyId(), errors.New("missing source id")
	}

	sourceId, err := types.IdFromString(rawSourceId)
	if err != nil {
		return types.EmptyId(), fmt.Errorf("malformed source id: %v", err)
	}

	return sourceId, nil
}

func getBearerToken(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("missing bearer token")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return parts[1], nil
}

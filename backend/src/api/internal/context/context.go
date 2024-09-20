package context

import (
	"errors"
	"fmt"
	"luna-backend/api/internal/config"
	"luna-backend/db"
	"luna-backend/types"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) *config.Api {
	return c.MustGet("apiConfig").(*config.Api)
}

func GetUserId(c *gin.Context) types.ID {
	return c.MustGet("user_id").(types.ID)
}

func GetId(c *gin.Context, primitive string) (types.ID, error) {
	rawId := c.Param(fmt.Sprintf("%sId", primitive))

	if rawId == "" {
		return types.EmptyId(), fmt.Errorf("missing %v id", primitive)
	}

	id, err := types.IdFromString(rawId)
	if err != nil {
		return types.EmptyId(), fmt.Errorf("malformed %v id: %v", primitive, err)
	}

	return id, nil
}

func GetBearerToken(c *gin.Context) (string, error) {
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

func GetTransaction(c *gin.Context) *db.Transaction {
	return c.MustGet("transaction").(*db.Transaction)
}

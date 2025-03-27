package util

import (
	"context"
	"fmt"
	"luna-backend/config"
	"luna-backend/db"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HandlerUtility struct {
	Config       *config.CommonConfig
	Logger       *logrus.Entry
	Tx           *db.Transaction
	Context      context.Context
	ResponseChan chan *Response
	ErrChan      chan *errors.ErrorTrace
	WarnChan     chan *errors.ErrorTrace
}

type Response struct {
	httpCode int
	msg      *gin.H
	file     types.File
	raw      []byte
	rawType  string
}

func (r *Response) GetStatus() int {
	return r.httpCode
}

func (r *Response) GetRaw() []byte {
	return r.raw
}

func (r *Response) GetRawType() string {
	return r.rawType
}

func (r *Response) GetMsg() *gin.H {
	return r.msg
}

func (r *Response) GetFile() types.File {
	return r.file
}

func (u *HandlerUtility) Success(msg *gin.H) {
	u.ResponseWithStatus(http.StatusOK, msg)
}

func (u *HandlerUtility) SuccessRawJson(rawJson []byte) {
	u.ResponseChan <- &Response{http.StatusOK, nil, nil, rawJson, "application/json"}
}

func (u *HandlerUtility) ResponseWithStatus(httpCode int, msg *gin.H) {
	u.ResponseChan <- &Response{httpCode, msg, nil, nil, ""}
}

func (u *HandlerUtility) ResponseWithFile(file types.File) {
	u.ResponseChan <- &Response{http.StatusOK, nil, file, nil, ""}
}

func (u *HandlerUtility) Error(err *errors.ErrorTrace) {
	u.ErrChan <- err
}

func (u *HandlerUtility) Warn(err *errors.ErrorTrace) {
	u.WarnChan <- err
}

func GetUtil(c *gin.Context) *HandlerUtility {
	return c.MustGet("handlerUtil").(*HandlerUtility)
}

func GetUserId(c *gin.Context) types.ID {
	return c.MustGet("user_id").(types.ID)
}

func GetId(c *gin.Context, primitive string) (types.ID, *errors.ErrorTrace) {
	rawId := c.Param(fmt.Sprintf("%sId", primitive))

	if rawId == "" {
		return types.EmptyId(), errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Missing %v id", primitive).
			Append(errors.LvlPlain, "Malformed request")
	}

	id, err := types.IdFromString(rawId)
	if err != nil {
		return types.EmptyId(), errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Malformed %v id", primitive).
			Append(errors.LvlPlain, "Malformed request")
	}

	return id, nil
}

func GetBearerToken(c *gin.Context) (string, *errors.ErrorTrace) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		return "", errors.New().Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "Missing authorization token").
			Append(errors.LvlPlain, "Unauthorized")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New().Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "Malformed authorization token").
			Append(errors.LvlPlain, "Unauthorized")
	}

	return parts[1], nil
}

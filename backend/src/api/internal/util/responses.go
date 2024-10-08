package util

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrorInvalidCredentials = iota
	ErrorDatabase           // TODO: a lot of these should be NotFound instead, but same problem as with ErrorUnknown
	ErrorInternal
	ErrorSourceNotFound
	ErrorCalendarNotFound
	ErrorEventNotFound
	ErrorPayload
	ErrorMalformedID
	ErrorTokenMissing
	ErrorTokenInvalid
	ErrorNotImplemented
	ErrorUnknown // TODO: all errors that use this cannot distinguish between types of errors; need to revamp error handling to know the source of error
)

var errorResponses = []*errorResponse{
	err(http.StatusUnauthorized, "Invalid credentials"),
	err(http.StatusInternalServerError, "Database error"),
	err(http.StatusInternalServerError, "Internal error"),
	err(http.StatusNotFound, "Source does not exist"),
	err(http.StatusNotFound, "Calendar does not exist"),
	err(http.StatusNotFound, "Event does not exist"),
	err(http.StatusBadRequest, "Malformed payload"),
	err(http.StatusBadRequest, "Malformed ID"),
	err(http.StatusUnauthorized, "Missing authorization token"),
	err(http.StatusUnauthorized, "Invalid authorization token"),
	err(http.StatusNotImplemented, "Not implemented"),
	err(http.StatusInternalServerError, "Unknown error"),
}

const (
	DetailName     = "Malformed or missing name"
	DetailAuth     = "Malformed or missing authentication"
	DetailPassword = "Malformed or missing password"
	DetailSource   = "Malformed or missing source parameters"
	DetailColor    = "Malformed or missing color"
	DetailId       = "Malformed or missing ID"
	DetailDate     = "Malformed or missing date"
	DetailFields   = "Nothing to update"
	DetailTime     = "Malformed or missing time"
)

type errorResponse struct {
	code int
	msg  string
}

func err(code int, msg string) *errorResponse {
	return &errorResponse{
		code: code,
		msg:  msg,
	}
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func Error(c *gin.Context, code int) {
	resp := errorResponses[code]
	c.JSON(resp.code, gin.H{"error": resp.msg})
}

func ErrorDetailed(c *gin.Context, code int, details string) {
	resp := errorResponses[code]
	c.JSON(resp.code, gin.H{"error": fmt.Sprintf("%s: %s", resp.msg, details)})
}

func Abort(c *gin.Context, code int) {
	resp := errorResponses[code]
	c.AbortWithStatusJSON(resp.code, gin.H{"error": resp.msg})
}

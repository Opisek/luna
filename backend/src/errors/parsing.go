package errors

import (
	"net/http"
	"strings"
)

func InterpretRemoteError(err error, resource string, wordyResource string) *ErrorTrace {
	errMsg := err.Error()

	switch {
	// No error
	case err == nil:
		return nil

	// HTTP response codes
	case strings.Contains(errMsg, http.StatusText(http.StatusUnauthorized)):
		return New().Status(http.StatusUnauthorized).
			Append(LvlWordy, "Unauthorized access to %v", wordyResource).
			AltStr(LvlPlain, "Wrong credentials")
	case strings.Contains(errMsg, http.StatusText(http.StatusNotFound)):
		return New().Status(http.StatusNotFound).
			Append(LvlWordy, "%v not found", wordyResource).
			AltStr(LvlPlain, "%v not found", resource)
	case strings.Contains(errMsg, http.StatusText(http.StatusForbidden)):
		return New().Status(http.StatusForbidden).
			Append(LvlWordy, "Access to %v forbidden", wordyResource).
			AltStr(LvlPlain, "Access forbidden")
	case strings.Contains(errMsg, http.StatusText(http.StatusServiceUnavailable)):
		return New().Status(http.StatusServiceUnavailable).
			Append(LvlWordy, "%v temporarily unavailable", wordyResource).
			AltStr(LvlPlain, "%v temporarily unavailable", resource)
	case strings.Contains(errMsg, http.StatusText(http.StatusInternalServerError)):
		return New().Status(http.StatusInternalServerError).
			Append(LvlPlain, "Remote server returned an error")
	case strings.Contains(errMsg, http.StatusText(http.StatusBadRequest)):
		fallthrough
	case strings.Contains(errMsg, http.StatusText(http.StatusMethodNotAllowed)):
		return New().Status(http.StatusBadRequest).
			Append(LvlWordy, "Bad request to %v", wordyResource).
			AltStr(LvlPlain, "Bad request to %v, are you sure the URL is correct?", resource)

	// Connection errors
	case strings.Contains(errMsg, "dial tcp"):
		return New().Status(http.StatusServiceUnavailable).
			Append(LvlWordy, "Could not connect to %v", wordyResource).
			AltStr(LvlPlain, "Could not connect to %v", resource)
	case strings.Contains(errMsg, "no such host"):
		return New().Status(http.StatusServiceUnavailable).
			Append(LvlWordy, "Could not resolve %v", wordyResource).
			AltStr(LvlPlain, "Could not resolve %v", resource)
	case strings.Contains(errMsg, "connection refused"):
		return New().Status(http.StatusServiceUnavailable).
			Append(LvlWordy, "Connection to %v refused", wordyResource).
			AltStr(LvlPlain, "Connection refused to %v", resource)
	case strings.Contains(errMsg, "connection reset by peer"):
		return New().Status(http.StatusServiceUnavailable).
			Append(LvlWordy, "Connection to %v reset by peer", wordyResource).
			AltStr(LvlPlain, "Connection to %v reset by peer", resource)

	// Other errors
	default:
		return New().Status(http.StatusInternalServerError).
			AddErr(LvlDebug, err).
			Append(LvlDebug, "Could not query CalDAV source")
	}
}

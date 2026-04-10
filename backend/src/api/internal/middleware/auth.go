package middleware

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/crypto"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
)

// Error messages must be kept intentionally vague even in higher verbosity levels,
// to avoid the creation of an oracle.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := util.GetUtil(c)

		cookieToken, cookieErr := c.Cookie("token")
		gotCookie := cookieErr == nil && cookieToken != ""
		bearerToken, bearerErr := util.GetBearerToken(c)
		gotBearer := bearerErr == nil && bearerToken != ""

		if !gotCookie && !gotBearer {
			u.Error(errors.New().Status(http.StatusUnauthorized).
				Append(errors.LvlWordy, "Missing token"))
			c.Abort()
			return
		}

		var token string
		if gotBearer {
			token = bearerToken
		} else if gotCookie {
			token = cookieToken
		}

		// Parse the token => Verifies that it was signed by the server
		parsedToken, tr := auth.ParseToken(u.Config, token)
		if tr != nil {
			u.Error(tr)
			c.Abort()
			return
		}

		// Find the session in the database => Verifies that the session has not been revoked and that the user exists
		session, tr := u.Tx.Queries().GetSessionAndUpdateLastSeen(parsedToken.UserId, parsedToken.SessionId, util.DetermineClientAddress(c))
		if tr != nil {
			u.Error(tr)
			c.Abort()
			return
		}

		// Check if the secret from the token hashes to the same value as the one in the database.
		// This is one line of defense against forged tokens.
		secret, err := base64.StdEncoding.DecodeString(parsedToken.Secret)
		if err != nil {
			u.Error(errors.New().Status(http.StatusUnauthorized).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not decode secret from token").
				Append(errors.LvlPlain, "Session expired"),
			)
			c.Abort()
			return
		}
		serverSecret, tr := crypto.GetSymmetricKey(u.Config, "tokenHashSecret")
		if tr != nil {
			u.Error(tr)
			c.Abort()
			return
		}
		actualHash := crypto.GetSha256Hash(serverSecret, parsedToken.SessionId.Bytes(), secret)
		if !bytes.Equal(actualHash, session.SecretHash) {
			u.Error(errors.New().Status(http.StatusUnauthorized).
				Append(errors.LvlDebug, "Token secret value produces incorrect hash value").
				Append(errors.LvlPlain, "Session expired"),
			)
			c.Abort()
			return
		}

		// Check if the user agent matches the one used to create the session
		// This is not a 100% guarantee that the token was not stolen, but it is a good first line of defense
		// This is disabled for API requests and for loading data from +page.ts or +layout.ts in Svelte.
		// While an attacker could just set their user agent to "FrontendLoad",
		// this is not that much easier than setting the correct user agent in first place.
		if !(session.IsApi || (c.Request.UserAgent() == "FrontendLoad" && (c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead))) {
			// Compare the user agents => Verifies that the token was not stolen and used from another device (not a 100% guarantee)
			associatedUserAgent := useragent.Parse(session.UserAgent)
			currentUserAgent := useragent.Parse(c.Request.UserAgent())

			osMatch := associatedUserAgent.OS == currentUserAgent.OS
			nameMatch := associatedUserAgent.Name == currentUserAgent.Name
			deviceMatch := associatedUserAgent.Device == currentUserAgent.Device

			if !osMatch || !nameMatch || !deviceMatch || (session.IsShortLived && associatedUserAgent.String != currentUserAgent.String) {
				u.Error(errors.New().Status(http.StatusUnauthorized).
					Append(errors.LvlDebug, "Expected %s, but got %s",
						fmt.Sprintf("%s %s %s", associatedUserAgent.OS, associatedUserAgent.Name, associatedUserAgent.Device),
						fmt.Sprintf("%s %s %s", currentUserAgent.OS, currentUserAgent.Name, currentUserAgent.Device)).
					Append(errors.LvlDebug, "User agent mismatch").
					Append(errors.LvlPlain, "Session expired"),
				)

				u.Config.TokenInvalidationChannel <- session

				c.Abort()
				return
			}
		}

		// Get the permissions associated with the token
		var permissions *types.TokenPermissions
		if !session.IsApi {
			permissions = types.AllPermissions()
		} else {
			permissions, tr = u.Tx.Queries().GetTokenPermissions(parsedToken.SessionId)
			if tr != nil {
				u.Error(tr)
				c.Abort()
				return
			}
		}

		c.Set("user_id", parsedToken.UserId)
		c.Set("session_id", parsedToken.SessionId)
		c.Set("permissions", permissions)

		c.Next()
	}
}

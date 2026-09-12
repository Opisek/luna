package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/auth"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// OAuth 2.0 Clients
//

func GetOauthClient(c *gin.Context) {
	u := util.GetUtil(c)

	// Client ID
	clientId, tr := util.GetId(c, "client")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Get the client
	client, tr := u.Tx.Queries().GetOauthClientById(clientId)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{
		"client": client,
	})
}

func GetOauthClients(c *gin.Context) {
	u := util.GetUtil(c)

	// Get the clients
	clients, tr := u.Tx.Queries().GetOauthClients()
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{
		"clients": clients,
	})
}

func PutOauthClient(c *gin.Context, body *struct {
	Name         string     `json:"name" form:"name" binding:"required,alphanumunicode"`
	ClientId     string     `json:"client_id" form:"client_id" binding:"required"`
	ClientSecret string     `json:"client_secret" form:"client_secret" binding:"required"`
	BaseUrl      *types.Url `json:"base_url" form:"base_url" binding:"required"`
	Scope        string     `json:"scope" form:"scope" binding:"required"`
}) {
	u := util.GetUtil(c)

	client := &types.OauthClient{
		Name:         body.Name,
		ClientId:     body.ClientId,
		ClientSecret: body.BaseUrl.Scheme,
		BaseUrl:      body.BaseUrl,
		Scope:        body.Scope,
	}

	// Check if the client is valid
	tr := auth.FetchOauthUrls(client, u.Context)
	if tr != nil {
		u.Error(tr.
			Append(errors.LvlWordy, "Could not determine OIDC URLs"),
		)
		return
	}

	// Insert
	tr = u.Tx.Queries().InsertOauthClient(client)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{
		"client": client,
	})
}

func PatchOauthClient(c *gin.Context, body *struct {
	Name         string     `json:"name" form:"name" binding:"required,alphanumunicode"`
	ClientId     string     `json:"client_id" form:"client_id" binding:"required"`
	ClientSecret string     `json:"client_secret" form:"client_secret" binding:"required"`
	BaseUrl      *types.Url `json:"base_url" form:"base_url" binding:"required"`
	Scope        string     `json:"scope" form:"scope" binding:"required"`
}) {
	u := util.GetUtil(c)

	// Client ID
	clientId, tr := util.GetId(c, "client")
	if tr != nil {
		u.Error(tr)
		return
	}

	client := &types.OauthClient{
		Id:           clientId,
		Name:         body.Name,
		ClientId:     body.ClientId,
		ClientSecret: body.ClientSecret,
		BaseUrl:      body.BaseUrl,
		Scope:        body.Scope,
	}

	// Check if the client is valid
	tr = auth.FetchOauthUrls(client, u.Context)
	if tr != nil {
		u.Error(tr.
			Append(errors.LvlWordy, "Could not determine OIDC URLs"),
		)
		return
	}

	// Update
	tr = u.Tx.Queries().UpdateOauthClient(client)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{
		"client": client,
	})
}

func DeleteOauthClient(c *gin.Context) {
	u := util.GetUtil(c)

	// Client ID
	clientId, tr := util.GetId(c, "client")
	if tr != nil {
		u.Error(tr)
		return
	}

	tr = u.Tx.Queries().DeleteOauthClient(clientId)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(nil)
}

//
// OAuth 2.0 Authorization Requests
//

func CreateOauthAuthorizationRequest(c *gin.Context) {
	u := util.GetUtil(c)

	// Client ID
	clientId, tr := util.GetId(c, "client")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Check if tokens already exist
	request := &types.OauthAuthorizationRequest{
		ClientId: clientId,
		UserId:   util.GetUserId(c),
	}

	// Insert
	tr = u.Tx.Queries().InsertOauthAuthorizationRequest(request)
	if tr != nil {
		u.Error(tr)
		return
	}

	// Build URL
	client, tr := u.Tx.Queries().GetOauthClientById(clientId)
	if tr != nil {
		u.Error(tr)
		return
	}

	tr = auth.FetchOauthUrls(client, u.Context)
	if tr != nil {
		u.Error(tr)
		return
	}

	consentUrl := *client.AuthorizationUrl.URL()
	queryParams := consentUrl.Query()

	// RFC 6749 4.1.1
	queryParams.Add("response_type", "code")
	queryParams.Add("client_id", client.ClientId)
	queryParams.Add("redirect_uri", auth.GetOauthRedirectUrl(u.Config).String())
	if client.Scope != "" {
		queryParams.Add("scope", client.Scope)
	}
	queryParams.Add("state", request.Id.String())

	// This is non-standard but required for Google: https://developers.google.com/identity/protocols/oauth2/web-server#offline
	// TODO: See if this disturbs other authorization providers, in which case we need to make this configurable.
	queryParams.Add("access_type", "offline")

	consentUrl.RawQuery = queryParams.Encode()

	u.Success(&gin.H{
		"request": request,
		"url":     consentUrl.String(),
	})
}

func FinalizeOauthAuthorizationRequest(c *gin.Context, body *struct {
	AuthCode string `json:"authorization_code" form:"authorization_code" binding:"required"`
}) {
	u := util.GetUtil(c)

	// Request ID
	requestId, tr := util.GetId(c, "request")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Fetch outstanding request
	request, tr := u.Tx.Queries().GetOauthAuthorizationRequest(requestId)
	if tr != nil {
		u.Error(tr)
		return
	}

	// If the user ID does not match, error
	if request.UserId != util.GetUserId(c) {
		u.Error(errors.New().Status(http.StatusUnauthorized).
			Append(errors.LvlDebug, "The outstanding request %v is not associated with the executing user %v", requestId, request.UserId).
			Append(errors.LvlWordy, "The outstanding request is not associated with the executing user"),
		)
		return
	}

	// Fetch OAuth 2.0 client
	client, tr := u.Tx.Queries().GetOauthClientById(request.ClientId)
	if tr != nil {
		u.Error(tr)
		return
	}

	// Use the authorization code to fetch tokens
	tr = auth.FetchOauthUrls(client, u.Context)
	if tr != nil {
		u.Error(tr)
		return
	}

	tokens, tr := auth.FetchOauthTokensUsingAuthorizationCode(client, body.AuthCode, u.Context, u.Config)
	if tr != nil {
		u.Error(tr)
		return
	}
	tokens.UserId = request.UserId

	// Get userinfo
	accountId, accountName, tr := auth.FetchOauthAccountIdentifier(client, request.UserId, tokens.AccessToken, u.Context)
	if tr != nil {
		u.Error(tr)
		return
	}
	tokens.AccountId = accountId
	tokens.AccountName = accountName

	// Check if tokens already exist
	exist, tr := u.Tx.Queries().CheckOauthTokensExist(request.ClientId, request.UserId, tokens.AccountId)
	if tr != nil {
		u.Error(tr)
		return
	}

	if exist {
		u.Warn(errors.New().
			Append(errors.LvlDebug, "OAuth 2.0 tokens for client %v (user %v, account %v) already exist", request.ClientId, request.UserId, tokens.AccountId).
			AltStr(errors.LvlWordy, "OAuth 2.0 tokens already exist").
			AltStr(errors.LvlPlain, "Already signed in"),
		)
	} else {
		// Save tokens
		tr = u.Tx.Queries().InsertOauthTokens(tokens)
		if tr != nil {
			u.Error(tr)
			return
		}
	}

	// Remove all matching requests
	tr = u.Tx.Queries().DeleteOauthAuthorizationRequests(request.ClientId, request.UserId)
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{
		"id": tokens.Id,
	})
}

func CancelOauthAuthorizationRequest(c *gin.Context) {
	u := util.GetUtil(c)

	// Request ID
	requestId, tr := util.GetId(c, "request")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Delete
	tr = u.Tx.Queries().DeleteOauthAuthorizationRequest(requestId, util.GetUserId(c))
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(nil)
}

//
// OAuth 2.0 Tokens
//

func GetOauthClientsWithTokens(c *gin.Context) {
	u := util.GetUtil(c)

	// Get IDs
	tokens, tr := u.Tx.Queries().GetOauthClientIdsWithTokens(util.GetUserId(c))
	if tr != nil {
		u.Error(tr)
		return
	}

	u.Success(&gin.H{
		"tokens": tokens,
	})
}

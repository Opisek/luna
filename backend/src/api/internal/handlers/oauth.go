package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"luna-backend/oauth"
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

func PutOauthClient(c *gin.Context) {
	u := util.GetUtil(c)

	// Oauth client name
	clientName := c.Request.FormValue("name")
	if clientName == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Name may not be empty"),
		)
		return
	}

	// Oauth client ID
	clientId := c.Request.FormValue("client_id")
	if clientId == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Client ID may not be empty"),
		)
		return
	}

	// Oauth client secret
	clientSecret := c.Request.FormValue("client_secret")
	if clientSecret == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Client secret may not be empty"),
		)
		return
	}

	// Oauth client base URL
	rawBaseUrl := c.Request.FormValue("base_url")
	if rawBaseUrl == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Base URL may not be empty"),
		)
		return
	}

	var err error
	err = util.IsValidUrl(rawBaseUrl)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Invalid base URL").
			AddErr(errors.LvlWordy, err),
		)
		return
	}
	baseUrl, err := types.NewUrl(rawBaseUrl)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Invalid base URL").
			AddErr(errors.LvlWordy, err),
		)
		return
	}

	// Check if the client is valid
	client := &types.OauthClient{
		Name:         clientName,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		BaseUrl:      baseUrl,
		Scope:        c.Request.FormValue("scope"),
	}

	tr := oauth.FetchOauthUrls(client, u.Context)
	if tr != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid base URL").
			AddErr(errors.LvlWordy, err),
		)
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

func PatchOauthClient(c *gin.Context) {
	u := util.GetUtil(c)

	// Client ID
	clientId, tr := util.GetId(c, "client")
	if tr != nil {
		u.Error(tr)
		return
	}

	// New client name
	newClientName := c.Request.FormValue("name")
	if newClientName == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Name may not be empty"),
		)
		return
	}

	// New client ID
	newClientId := c.Request.FormValue("client_id")
	if newClientId == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Client ID may not be empty"),
		)
		return
	}

	// New client secret
	newClientSecret := c.Request.FormValue("client_secret")

	// New client base URL
	rawBaseUrl := c.Request.FormValue("base_url")
	if rawBaseUrl == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Base URL may not be empty"),
		)
		return
	}

	var err error
	err = util.IsValidUrl(rawBaseUrl)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Invalid base URL").
			AddErr(errors.LvlWordy, err),
		)
		return
	}

	newBaseUrl, err := types.NewUrl(rawBaseUrl)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "Invalid base URL").
			AddErr(errors.LvlWordy, err),
		)
		return
	}

	// Check if the client is valid
	client := &types.OauthClient{
		Id:           clientId,
		Name:         newClientName,
		ClientId:     newClientId,
		ClientSecret: newClientSecret,
		BaseUrl:      newBaseUrl,
		Scope:        c.Request.FormValue("scope"),
	}

	tr = oauth.FetchOauthUrls(client, u.Context)
	if tr != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid base URL").
			AddErr(errors.LvlWordy, err),
		)
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

	exist, tr := u.Tx.Queries().CheckOauthTokensExist(request.ClientId, request.UserId)
	if tr != nil {
		u.Error(tr)
		return
	}
	if exist {
		u.Error(errors.New().Status(http.StatusConflict).
			Append(errors.LvlDebug, "OAuth 2.0 tokens for client %v (user %v) already exist", request.ClientId, request.UserId).
			AltStr(errors.LvlWordy, "OAuth 2.0 tokens already exist").
			AltStr(errors.LvlPlain, "Already signed in"),
		)
		return
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

	tr = oauth.FetchOauthUrls(client, c)
	if tr != nil {
		u.Error(tr)
		return
	}

	consentUrl := *client.AuthorizationUrl.URL()
	queryParams := consentUrl.Query()

	// RFC 6749 4.1.1
	queryParams.Add("response_type", "code")
	queryParams.Add("client_id", client.ClientId)
	queryParams.Add("redirect_uri", oauth.GetOauthRedirectUrl(u.Config).String())
	if client.Scope != "" {
		queryParams.Add("scope", client.Scope)
	}
	queryParams.Add("state", request.Id.String())

	consentUrl.RawQuery = queryParams.Encode()

	u.Success(&gin.H{
		"request": request,
		"url":     consentUrl.String(),
	})
}

func FinalizeOauthAuthorizationRequest(c *gin.Context) {
	u := util.GetUtil(c)

	// Request ID
	requestId, tr := util.GetId(c, "request")
	if tr != nil {
		u.Error(tr)
		return
	}

	// Authorization code
	authCode := c.Request.FormValue("authorization_code")
	if authCode == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlWordy, "Missing authorization coder"),
		)
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
	tr = oauth.FetchOauthUrls(client, c)
	if tr != nil {
		u.Error(tr)
		return
	}

	tokens, tr := oauth.FetchOauthTokensUsingAuthorizationCode(client, authCode, c, u.Config)
	if tr != nil {
		u.Error(tr)
		return
	}
	tokens.UserId = request.UserId

	// Check if tokens already exist
	exist, tr := u.Tx.Queries().CheckOauthTokensExist(request.ClientId, request.UserId)
	if tr != nil {
		u.Error(tr)
		return
	}

	if exist {
		u.Warn(errors.New().
			Append(errors.LvlDebug, "OAuth 2.0 tokens for client %v (user %v) already exist", request.ClientId, request.UserId).
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

	u.Success(nil)
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

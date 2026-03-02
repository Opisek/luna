package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	// Oauth client authorization URL
	rawAuthUrl := c.Request.FormValue("authorization_url")
	if rawAuthUrl == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Authorization URL may not be empty"),
		)
		return
	}

	var err error
	err = util.IsValidUrl(rawAuthUrl)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid authorization URL").
			AddErr(errors.LvlWordy, err),
		)
		return
	}
	authUrl, err := types.NewUrl(rawAuthUrl)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid authorization URL").
			AddErr(errors.LvlWordy, err),
		)
		return
	}

	// Insert
	client := &types.OauthClient{
		Name:             clientName,
		ClientId:         clientId,
		ClientSecret:     clientSecret,
		AuthorizationUrl: authUrl,
	}

	tr := u.Tx.Queries().InsertOauthClient(client)
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

	// New client authorization URL
	rawAuthUrl := c.Request.FormValue("authorization_url")
	if rawAuthUrl == "" {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Authorization URL may not be empty"),
		)
		return
	}

	var err error
	err = util.IsValidUrl(rawAuthUrl)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid authorization URL").
			AddErr(errors.LvlWordy, err),
		)
		return
	}

	newAuthUrl, err := types.NewUrl(rawAuthUrl)
	if err != nil {
		u.Error(errors.New().Status(http.StatusBadRequest).
			Append(errors.LvlPlain, "Invalid authorization URL").
			AddErr(errors.LvlWordy, err),
		)
		return
	}

	// Update
	client := &types.OauthClient{
		Id:               clientId,
		Name:             newClientName,
		ClientId:         newClientId,
		ClientSecret:     newClientSecret,
		AuthorizationUrl: newAuthUrl,
	}

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

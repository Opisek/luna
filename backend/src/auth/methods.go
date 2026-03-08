package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"luna-backend/config"
	"luna-backend/constants"
	"luna-backend/db"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"time"
)

// Wrapper

type httpClient struct {
	auth types.AuthMethod
}

func (httpClient *httpClient) Do(req *http.Request) (*http.Response, error) {
	res, tr := httpClient.auth.Do(req)
	if tr != nil {
		return nil, tr.SerializeError(errors.LvlDebug)
	}
	return res, nil
}

// No Authentication

type NoAuth struct{}

func (auth NoAuth) Do(req *http.Request) (*http.Response, *errors.ErrorTrace) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New().AddErr(errors.LvlDebug, err)
	}
	return res, nil
}

func (auth NoAuth) GetType() string {
	return constants.AuthNone
}

func (auth NoAuth) String() (string, error) {
	return "", nil
}

func (auth NoAuth) HttpClient() types.HttpClientInterface {
	return &httpClient{auth: auth}
}

func NewNoAuth() types.AuthMethod {
	return NoAuth{}
}

// Password and Username

type BasicAuth struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (auth BasicAuth) Do(req *http.Request) (*http.Response, *errors.ErrorTrace) {
	req.SetBasicAuth(auth.Username, auth.Password)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New().AddErr(errors.LvlDebug, err)
	}
	return res, nil
}

func (auth BasicAuth) GetType() string {
	return constants.AuthBasic
}

func (auth BasicAuth) String() (string, error) {
	bytes, err := json.Marshal(auth)
	if err != nil {
		return "", fmt.Errorf("could not marshal basic auth: %v", err)
	}
	return string(bytes), nil
}

func (auth BasicAuth) HttpClient() types.HttpClientInterface {
	return &httpClient{auth: auth}
}

func NewBasicAuth(username, password string) types.AuthMethod {
	return BasicAuth{Username: username, Password: password}
}

// Bearer Token

type BearerAuth struct {
	Token string `json:"token" form:"token"`
}

func (auth BearerAuth) Do(req *http.Request) (*http.Response, *errors.ErrorTrace) {
	req.Header.Set("Authorization", "Bearer "+auth.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New().AddErr(errors.LvlDebug, err)
	}
	return res, nil
}

func (auth BearerAuth) GetType() string {
	return constants.AuthBearer
}

func (auth BearerAuth) String() (string, error) {
	bytes, err := json.Marshal(auth)
	if err != nil {
		return "", fmt.Errorf("could not marshal bearer auth: %v", err)
	}
	return string(bytes), nil
}

func (auth BearerAuth) HttpClient() types.HttpClientInterface {
	return &httpClient{auth: auth}
}

func NewBearerAuth(token string) types.AuthMethod {
	return BearerAuth{Token: token}
}

// OAuth2

type OauthAuth struct {
	ClientId types.ID             `json:"client_id" form:"client_id"`
	client   *types.OauthClient   `json:"-" form:""`
	config   *config.CommonConfig `json:"-" form:""`
	ctx      context.Context      `json:"-" form:""`
	tx       *db.Transaction      `json:"-" form:""`
	userId   types.ID             `json:"-" form:""`
}

func (auth *OauthAuth) SupplyContext(userId types.ID, ctx context.Context, config *config.CommonConfig) {
	auth.config = config
	auth.ctx = ctx
	auth.userId = userId

	// I really don't like having to do the following, but I found no other way to
	// bring the current transaction into this struct without running into a lot of circular dependencies.
	// Sadly, I did not think this far when initially designing the auth package, because OAuth 2.0 was a faraway future back then.
	// Passing this reference via the context at least limits the amount of unsafe casts that we have to do to one.
	// I would like to revisit this later, but if I cannot find a better way to go about this,
	// I must at least handle the "do not use strings as keys" situation.
	auth.tx = ctx.Value("transaction").(*db.Transaction)
}

func (auth *OauthAuth) Do(req *http.Request) (*http.Response, *errors.ErrorTrace) {
	var tr *errors.ErrorTrace

	if auth.client == nil {
		auth.client, tr = auth.tx.Queries().GetOauthClientById(auth.ClientId)
		if tr != nil {
			return nil, tr
		}
	}

	tokens, tr := auth.tx.Queries().GetOauthTokens(auth.ClientId, auth.userId)
	if tr != nil {
		return nil, tr
	}

	if tokens.Expires.Before(time.Now()) {
		tr = FetchOauthUrls(auth.client, auth.ctx)
		if tr != nil {
			return nil, tr
		}

		refreshToken := tokens.RefreshToken
		tokens, tr = FetchOauthTokensUsingRefreshToken(auth.client, refreshToken, auth.ctx, auth.config)
		if tr != nil {
			return nil, tr
		}

		if tokens.RefreshToken == "" {
			tokens.RefreshToken = refreshToken
		}

		tr = auth.tx.Queries().UpdateOauthTokens(tokens)
		if tr != nil {
			return nil, tr
		}
	}

	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New().AddErr(errors.LvlDebug, err)
	}
	return res, nil
}

func (auth *OauthAuth) GetType() string {
	return constants.AuthOauth
}

func (auth *OauthAuth) String() (string, error) {
	bytes, err := json.Marshal(auth)
	if err != nil {
		return "", fmt.Errorf("could not marshal OAuth 2.0 auth: %v", err)
	}
	return string(bytes), nil
}

func (auth *OauthAuth) HttpClient() types.HttpClientInterface {
	return &httpClient{auth: auth}
}

func NewOauthAuth(clientId types.ID) types.AuthMethod {
	return &OauthAuth{ClientId: clientId}
}

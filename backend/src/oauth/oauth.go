package oauth

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"luna-backend/auth"
	"luna-backend/config"
	"luna-backend/errors"
	"luna-backend/net"
	"luna-backend/types"
)

type oauthUrlsResponse struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
}

func FetchOauthUrls(oauthClient *types.OauthClient, ctx context.Context) *errors.ErrorTrace {
	oidcUrl := oauthClient.BaseUrl.Subpage(".well-known", "openid-configuration")

	res := &oauthUrlsResponse{}

	tr := net.FetchJson(oidcUrl, "GET", auth.NewNoAuth(), nil, "", ctx, res)
	if tr != nil {
		return tr.
			Append(errors.LvlDebug, "Could not resolve OpenID connect configuration %v", oidcUrl.String()).
			Append(errors.LvlDebug, "Could not fetch URLs for OAuth 2.0 client %v", oauthClient.Name).
			AltStr(errors.LvlWordy, "Could not fetch URLs for OAuth 2.0 client")
	}

	var err error

	oauthClient.AuthorizationUrl, err = types.NewUrl(res.AuthorizationEndpoint)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse authorization URL %v", res.AuthorizationEndpoint).
			Append(errors.LvlDebug, "Could not fetch URLs for OAuth 2.0 client %v", oauthClient.Name).
			AltStr(errors.LvlWordy, "Could not fetch URLs for OAuth 2.0 client")
	}

	oauthClient.TokenUrl, err = types.NewUrl(res.TokenEndpoint)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse token URL %v", res.TokenEndpoint).
			Append(errors.LvlDebug, "Could not fetch URLs for OAuth 2.0 client %v", oauthClient.Name).
			AltStr(errors.LvlWordy, "Could not fetch URLs for OAuth 2.0 client")
	}

	return nil
}

func GetOauthRedirectUrl(config *config.CommonConfig) *types.Url {
	return config.PublicUrl.Subpage("/oauth")
}

// RFC 6749 5.1
type oauthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Expires      int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func FetchOauthTokensUsingAuthorizationCode(oauthClient *types.OauthClient, authCode string, ctx context.Context, config *config.CommonConfig) (*types.OauthTokens, *errors.ErrorTrace) {
	form := make(url.Values)

	// RFC 6749 4.1.3
	form.Add("grant_type", "authorization_code")
	form.Add("code", authCode)
	form.Add("redirect_uri", GetOauthRedirectUrl(config).String())
	form.Add("client_id", oauthClient.ClientId)
	form.Add("client_secret", oauthClient.ClientSecret)

	res := &oauthTokenResponse{}

	timestamp := time.Now()

	tr := net.FetchJson(oauthClient.TokenUrl, "POST", auth.NewNoAuth(), &form, "application/x-www-form-urlencoded", ctx, res)
	if tr != nil {
		return nil, tr.
			Append(errors.LvlDebug, "Could not fetch tokens for OAuth 2.0 client %v", oauthClient.Name).
			AltStr(errors.LvlWordy, "Could not fetch tokens for OAuth 2.0 client")
	}

	if res.Scope != "" && res.Scope != oauthClient.Scope {
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Returned scope %v does not match the requested scope %v", res.Scope, oauthClient.Scope).
			Append(errors.LvlDebug, "Could not fetch tokens for OAuth 2.0 client %v", oauthClient.Name).
			AltStr(errors.LvlWordy, "Could not fetch tokens for OAuth 2.0 client")
	}

	if res.Expires == 0 {
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Omitting explicit token expiry duration is not currently supported").
			Append(errors.LvlDebug, "Could not fetch tokens for OAuth 2.0 client %v", oauthClient.Name).
			AltStr(errors.LvlWordy, "Could not fetch tokens for OAuth 2.0 client")
	}

	if res.AccessToken == "" {
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Received an empty access token").
			Append(errors.LvlDebug, "Could not fetch tokens for OAuth 2.0 client %v", oauthClient.Name).
			AltStr(errors.LvlWordy, "Could not fetch tokens for OAuth 2.0 client")
	}

	timestamp = timestamp.Add(time.Duration(res.Expires) * time.Second)

	return &types.OauthTokens{
		ClientId:     oauthClient.Id,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		Expires:      timestamp,
	}, nil
}

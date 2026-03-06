package oauth

import (
	"context"
	"net/http"

	"luna-backend/auth"
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

	tr := net.FetchJson(oidcUrl, auth.NewNoAuth(), ctx, res)
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

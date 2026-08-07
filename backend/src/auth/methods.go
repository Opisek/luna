package auth

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"luna-backend/config"
	"luna-backend/constants"
	"luna-backend/db"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
	"regexp"
	"strings"
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

type digestChallenge struct {
	Realm     string
	Nonce     string
	Opaque    string
	Qop       string
	Algorithm string
	Stale     string
}

var digestAuthHeaderPrefix = regexp.MustCompile(`(?i)^digest\s+`)

func md5Hex(value string) string {
	sum := md5.Sum([]byte(value))
	return hex.EncodeToString(sum[:])
}

func splitHeaderFields(header string) []string {
	parts := []string{}
	start := 0
	inQuotes := false

	for i := 0; i < len(header); i++ {
		switch header[i] {
		case '"':
			inQuotes = !inQuotes
		case ',':
			if !inQuotes {
				parts = append(parts, strings.TrimSpace(header[start:i]))
				start = i + 1
			}
		}
	}

	parts = append(parts, strings.TrimSpace(header[start:]))
	return parts
}

func parseDigestChallenge(header string) *digestChallenge {
	header = strings.TrimSpace(digestAuthHeaderPrefix.ReplaceAllString(header, ""))
	if header == "" {
		return nil
	}

	challenge := &digestChallenge{}
	for _, field := range splitHeaderFields(header) {
		if field == "" {
			continue
		}

		key, value, ok := strings.Cut(field, "=")
		if !ok {
			continue
		}

		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.Trim(strings.TrimSpace(value), `"`)

		switch key {
		case "realm":
			challenge.Realm = value
		case "nonce":
			challenge.Nonce = value
		case "opaque":
			challenge.Opaque = value
		case "qop":
			challenge.Qop = value
		case "algorithm":
			challenge.Algorithm = value
		case "stale":
			challenge.Stale = value
		}
	}

	if challenge.Realm == "" || challenge.Nonce == "" {
		return nil
	}

	return challenge
}

func cloneRequest(req *http.Request) (*http.Request, error) {
	cloned := req.Clone(req.Context())
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		cloned.Body = body
	} else if req.Body != nil && req.Body != http.NoBody {
		return nil, fmt.Errorf("request body is not replayable")
	} else {
		cloned.Body = nil
	}
	return cloned, nil
}

func newCnonce() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func digestAuthorization(req *http.Request, challenge *digestChallenge, username string, password string) (string, error) {
	algorithm := strings.ToUpper(challenge.Algorithm)
	if algorithm == "" {
		algorithm = "MD5"
	}
	if algorithm != "MD5" {
		return "", fmt.Errorf("unsupported digest algorithm %q", challenge.Algorithm)
	}

	qop := ""
	if challenge.Qop != "" {
		for _, option := range splitHeaderFields(challenge.Qop) {
			if strings.EqualFold(option, "auth") {
				qop = "auth"
				break
			}
		}
		if qop == "" {
			return "", fmt.Errorf("unsupported digest qop %q", challenge.Qop)
		}
	}

	cnonce, err := newCnonce()
	if err != nil {
		return "", err
	}

	uri := req.URL.RequestURI()
	ha1 := md5Hex(fmt.Sprintf("%s:%s:%s", username, challenge.Realm, password))
	ha2 := md5Hex(fmt.Sprintf("%s:%s", req.Method, uri))

	response := ""
	parts := []string{
		fmt.Sprintf(`username="%s"`, username),
		fmt.Sprintf(`realm="%s"`, challenge.Realm),
		fmt.Sprintf(`nonce="%s"`, challenge.Nonce),
		fmt.Sprintf(`uri="%s"`, uri),
	}

	if qop != "" {
		nonceCount := "00000001"
		response = md5Hex(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, challenge.Nonce, nonceCount, cnonce, qop, ha2))
		parts = append(parts,
			fmt.Sprintf(`qop=%s`, qop),
			fmt.Sprintf(`nc=%s`, nonceCount),
			fmt.Sprintf(`cnonce="%s"`, cnonce),
		)
	} else {
		response = md5Hex(fmt.Sprintf("%s:%s:%s", ha1, challenge.Nonce, ha2))
	}

	parts = append(parts, fmt.Sprintf(`response="%s"`, response))
	if challenge.Opaque != "" {
		parts = append(parts, fmt.Sprintf(`opaque="%s"`, challenge.Opaque))
	}
	if challenge.Algorithm != "" {
		parts = append(parts, fmt.Sprintf(`algorithm=%s`, challenge.Algorithm))
	}

	return "Digest " + strings.Join(parts, ", "), nil
}

func (auth BasicAuth) doDigest(req *http.Request, challenge *digestChallenge) (*http.Response, *errors.ErrorTrace) {
	retryReq, err := cloneRequest(req)
	if err != nil {
		return nil, errors.New().AddErr(errors.LvlDebug, err)
	}

	authorization, err := digestAuthorization(retryReq, challenge, auth.Username, auth.Password)
	if err != nil {
		return nil, errors.New().AddErr(errors.LvlDebug, err)
	}

	retryReq.Header.Set("Authorization", authorization)
	res, err := http.DefaultClient.Do(retryReq)
	if err != nil {
		return nil, errors.New().AddErr(errors.LvlDebug, err)
	}
	return res, nil
}

func (auth BasicAuth) Do(req *http.Request) (*http.Response, *errors.ErrorTrace) {
	basicReq, err := cloneRequest(req)
	if err != nil {
		return nil, errors.New().AddErr(errors.LvlDebug, err)
	}

	basicReq.SetBasicAuth(auth.Username, auth.Password)
	res, err := http.DefaultClient.Do(basicReq)
	if err != nil {
		return nil, errors.New().AddErr(errors.LvlDebug, err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		challenge := parseDigestChallenge(res.Header.Get("WWW-Authenticate"))
		if challenge != nil {
			io.Copy(io.Discard, res.Body)
			res.Body.Close()
			return auth.doDigest(req, challenge)
		}
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
	TokensID types.ID `json:"tokens_id" form:"tokens_id"`
	ClientId types.ID `json:"client_id" form:"client_id"`
	UserId   types.ID `json:"user_id" form:"user_id"`

	client *types.OauthClient `json:"-" form:""`
	ctx    context.Context    `json:"-" form:""`
	tx     *db.Transaction    `json:"-" form:""`
}

func (auth *OauthAuth) SupplyContext(ctx context.Context) {
	auth.ctx = ctx

	// I really don't like having to do the following hack, but I found no other way to
	// bring the current transaction into this struct without running into a lot of circular dependencies.
	// Sadly, I did not think this far when initially designing the auth package, because OAuth 2.0 was a faraway future back then.
	// Passing this reference via the context at least limits the amount of unsafe casts that we have to do to one.
	// I would like to revisit this later, but if I cannot find a better way to go about this,
	// I must at least handle the "do not use strings as keys" situation.
	auth.tx = ctx.Value("transaction").(*db.Transaction)
}

func (auth *OauthAuth) expired() *errors.ErrorTrace {
	// Same as above
	(auth.ctx.Value("config").(*config.CommonConfig)).OauthInvalidationChannel <- auth.TokensID

	return errors.New().Status(http.StatusUnauthorized).
		Append(errors.LvlDebug, "OAuth 2.0 tokens %v for client %v and user %v expired or were revoked", auth.TokensID, auth.ClientId, auth.UserId).
		AltStr(errors.LvlWordy, "OAuth 2.0 tokens expired or were revoked").
		Append(errors.LvlPlain, "Please authorize yourself again")
}

func (auth *OauthAuth) Do(req *http.Request) (*http.Response, *errors.ErrorTrace) {
	var tr *errors.ErrorTrace

	if auth.client == nil {
		auth.client, tr = auth.tx.Queries().GetOauthClientById(auth.ClientId)
		if tr != nil {
			return nil, tr
		}
	}

	tokens, tr := auth.tx.Queries().GetOauthTokens(auth.TokensID, auth.UserId)
	if tr != nil {
		return nil, tr
	}

	if tokens.Expires.Before(time.Now()) {
		tr = FetchOauthUrls(auth.client, auth.ctx)
		if tr != nil {
			return nil, tr
		}

		newTokens, tr := FetchOauthTokensUsingRefreshToken(auth.client, tokens.RefreshToken, auth.ctx)
		if tr != nil {
			if strings.Contains(tr.Serialize(errors.LvlDebug), "invalid_grant") {
				return nil, auth.expired()
			}
			return nil, tr
		}
		newTokens.Id = tokens.Id
		newTokens.UserId = tokens.UserId

		if newTokens.RefreshToken == "" {
			newTokens.RefreshToken = tokens.RefreshToken
		}

		tr = auth.tx.Queries().UpdateOauthTokens(newTokens)
		if tr != nil {
			return nil, tr
		}

		tokens = newTokens
	}

	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid_grant") {
			return nil, auth.expired()
		}
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

func NewOauthAuth(tokensId types.ID, clientId types.ID, userId types.ID) types.AuthMethod {
	return &OauthAuth{
		TokensID: tokensId,
		ClientId: clientId,
		UserId:   userId,
	}
}

package types

import "time"

type OauthClient struct {
	Id               ID     `json:"id" db:"id" encrypted:"false"`
	Name             string `json:"name" db:"name" encrypted:"false"`
	ClientId         string `json:"client_id" db:"client_id" encrypted:"false"`
	ClientSecret     string `json:"client_secret" db:"client_secret" encrypted:"true"`
	BaseUrl          *Url   `json:"base_url" db:"base_url" encrypted:"false"`
	AuthorizationUrl *Url   `json:"-" db:"" encrypted:"false"`
	TokenUrl         *Url   `json:"-" db:"" encrypted:"false"`
	Scope            string `json:"scope" db:"scope" encrypted:"false"`
}

type OauthAuthorizationRequest struct {
	Id       ID `json:"request_id" db:"request_id" encrypted:"false"`
	ClientId ID `json:"client_id" db:"client_id" encrypted:"false"`
	UserId   ID `json:"user_id" db:"user_id" encrypted:"false"`
}

type OauthTokens struct {
	ClientId     ID        `json:"client_id" db:"client_id" encrypted:"false"`
	UserId       ID        `json:"user_id" db:"user_id" encrypted:"false"`
	AccessToken  []byte    `json:"-" db:"access_token" encrypted:"true"`
	RefreshToken []byte    `json:"-" db:"refresh_token" encrypted:"true"`
	Expires      time.Time `json:"expires_at" db:"expires_at" encrypted:"false"`
}

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
	RevocationUrl    *Url   `json:"-" db:"" encrypted:"false"`
	UserinfoUrl      *Url   `json:"-" db:"" encrypted:"false"`
	Scope            string `json:"scope" db:"scope" encrypted:"false"`
}

type OauthAuthorizationRequest struct {
	Id       ID        `json:"request_id" db:"request_id" encrypted:"false"`
	ClientId ID        `json:"client_id" db:"client_id" encrypted:"false"`
	UserId   ID        `json:"user_id" db:"user_id" encrypted:"false"`
	Expires  time.Time `json:"expires_at" db:"expires_at" encrpted:"false"`
}

type OauthTokens struct {
	Id           ID        `json:"id" db:"id" encrypted:"false"`
	ClientId     ID        `json:"client_id" db:"client_id" encrypted:"false"`
	UserId       ID        `json:"-" db:"user_id" encrypted:"false"`
	AccountId    string    `json:"account_id" db:"account_id" encrypted:"false"`
	AccountName  string    `json:"account_name" db:"account_name" encrypted:"false"`
	AccessToken  string    `json:"-" db:"access_token" encrypted:"true"`
	RefreshToken string    `json:"-" db:"refresh_token" encrypted:"true"`
	Expires      time.Time `json:"-" db:"expires_at" encrypted:"false"`
}

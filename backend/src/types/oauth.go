package types

type OauthClient struct {
	Id               ID     `json:"id" db:"id" encrypted:"false"`
	Name             string `json:"name" db:"name" encrypted:"false"`
	ClientId         string `json:"client_id" db:"client_id" encrypted:"false"`
	ClientSecret     string `json:"client_secret" db:"client_secret" encrypted:"true"`
	AuthorizationUrl *Url   `json:"authorization_url" db:"authorization_url" encrypted:"false"`
}

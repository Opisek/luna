package types

type OauthClient struct {
	Id               ID     `json:"id" db:"id" encrypted:"false"`
	Name             string `json:"name" db:"name" encrypted:"false"`
	ClientId         string `json:"client_id" db:"client_id" encrypted:"false"`
	ClientSecret     string `json:"client_secret" db:"client_secret" encrypted:"true"`
	BaseUrl          *Url   `json:"base_url" db:"base_url" encrypted:"false"`
	AuthorizationUrl *Url   `json:"" db:"" encrypted:"false"`
	TokenUrl         *Url   `json:"" db:"" encrypted:"false"`
	Scope            string `json:"scope" db:"scope" encrypted:"false"`
}

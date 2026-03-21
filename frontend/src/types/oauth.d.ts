type OauthClientModel = {
  id: string;
  name: string;
  client_id: string;
  client_secret: string;
  base_url: string;
  scope: string;
}

type OauthTokensModel = {
  id: string;
  client_id: string;
  account_id: string;
  account_name: string;
}
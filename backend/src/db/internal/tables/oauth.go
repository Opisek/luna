package tables

func (q *Tables) InitializeOauthClientsTable() error {
	// Oauth tokens table:
	// id name client_id client_secret base_url scope
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE oauth_clients (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

			name VARCHAR(255) NOT NULL UNIQUE,
			client_id VARCHAR(1024) NOT NULL,
			client_secret BYTEA NOT NULL,
			base_url VARCHAR(2048) NOT NULL,
			scope VARCHAR(1024) NOT NULL
		);
		`,
	)

	return err
}

// Warning: The client_id in the following tables is NOT to be confused with the client_id in the clients table.
// When other tables refer to an OAuth 2.0 client, they do so by the internal UUID.
func (q *Tables) InitializeOauthAuthorizationRequestsTable() error {
	// Oauth authorization requests table:
	// request_id oauth_client_id user_id expires_at
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE oauth_authorization_requests (
			request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			client_id UUID NOT NULL REFERENCES oauth_clients(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			expires_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '10 minutes'
		)
		`,
	)

	return err
}

func (q *Tables) InitializeOauthTokensTable() error {
	// Oauth tokens table:
	// oauth_client_id user_id access_token refresh_token expires_at
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE oauth_tokens (
			client_id UUID NOT NULL REFERENCES oauth_clients(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			access_token BYTEA NOT NULL,
			refresh_token BYTEA,
			expires_at TIMESTAMP NOT NULL,
			PRIMARY KEY (client_id, user_id)
		);
		`,
	)

	return err
}

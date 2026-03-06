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

func (q *Tables) InitializeOauthTokensTable() error {
	// Oauth tokens table:
	// oath_client_id user_id access_token refresh_token expires_at
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE oauth_tokens (
			oauth_client_id UUID NOT NULL REFERENCES oauth_clients(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			access_token BYTEA NOT NULL,
			refresh_token BYTEA NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			PRIMARY KEY (oauth_client_id, user_id)
		);
		`,
	)

	return err
}

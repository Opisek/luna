package tables

func (q *Tables) InitializeOauthClientsTable() error {
	// Oauth table:
	// id name client_id client_secret authorization_url
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE oauth_clients (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

			name VARCHAR(255) NOT NULL UNIQUE,
			client_id VARCHAR(1024) NOT NULL,
			client_secret BYTEA NOT NULL,
			authorization_url VARCHAR(2048) NOT NULL
		);
		`,
	)

	return err
}

package tables

func (q *Tables) InitializeSessionsTable() error {
	// Sessions table:
	// sessionid userid created_at last_seen user_agent ip_address is_short_lived is_api hash
	//
	// The hash is an additional security measure:
	// When creating a token, a random secret is generated and stored in the
	// encrypted token. A hash of that secret is stored in the database and
	// compared with a freshly generated hash whenever the token is used. This
	// ensures that an attacker with read-only access to the database and the
	// cryptographic keys cannot forge or create valid tokens.
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE sessions (
			sessionid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			userid UUID REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			user_agent TEXT,
			initial_ip_address INET,
			last_ip_address INET,
			is_short_lived BOOLEAN DEFAULT FALSE,
			is_api BOOLEAN DEFAULT FALSE,
			hash BYTEA NOT NULL
		);
		`,
	)

	return err
}

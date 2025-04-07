package tables

func (q *Tables) InitializeSessionsTable() error {
	// Sessions table:
	// sessionid userid created_at last_seen user_agent ip_address is_short_lived is_api
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE sessions (
			sessionid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			userid UUID REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			user_agent TEXT,
			ip_address INET,
			is_short_lived BOOLEAN DEFAULT FALSE,
			is_api BOOLEAN DEFAULT FALSE
		);
		`,
	)

	return err
}

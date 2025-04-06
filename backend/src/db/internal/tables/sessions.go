package tables

func (q *Tables) InitializeSessionsTable() error {
	// Sessions table:
	// sessionid userid created_at last_seen user_agent is_api
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS passwords (
			sessionid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			userid UUID REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
			user_agent TEXT,
			is_api BOOLEAN DEFAULT FALSE,
		);
		`,
	)

	return err
}

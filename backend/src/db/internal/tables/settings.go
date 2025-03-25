package tables

func (q *Tables) InitializeUserSettingsTable() error {
	// User settings table:
	// userid key value
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS user_settings (
			userid UUID REFERENCES users(id) ON DELETE CASCADE,
			key VARCHAR(64) NOT NULL,
			value BYTEA NOT NULL,
			PRIMARY KEY (userid, key)
		);
		`,
	)

	return err
}

func (q *Tables) InitializeGlobalSettingsTable() error {
	// Global settings table:
	// key value
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS global_settings (
			key VARCHAR(64) PRIMARY KEY,
			value BYTEA NOT NULL
		);
		`,
	)

	return err
}

package tables

func (q *Tables) InitializeTokenPermissionsTable() error {
	// Token permissions table:
	// sessionid permission
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE token_permissions (
			sessionid UUID REFERENCES sessions(sessionid) ON DELETE CASCADE,
			permission VARCHAR(64) NOT NULL
		);
		`,
	)

	return err
}

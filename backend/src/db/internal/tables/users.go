package tables

func (q *Tables) InitializeUsersTable() error {
	// Auth table:
	// id username password email admin
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			admin BOOLEAN,
			searchable BOOLEAN,
			profile_picture VARCHAR(255)
		);
		`,
	)

	return err
}

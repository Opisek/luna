package tables

import "context"

func (q *Tables) InitializeUserTable() error {
	// Auth table:
	// id username password email admin
	_, err := q.Tx.Exec(
		context.TODO(),
		`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			algorithm VARCHAR(32) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			admin BOOLEAN
		);
		`,
	)

	return err
}

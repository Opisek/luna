package tables

import "context"

func (q *Tables) InitializePasswordsTable() error {
	// Auth table:
	// id hash salt algorithm parameters
	_, err := q.Tx.Exec(
		context.TODO(),
		`
		CREATE TABLE IF NOT EXISTS passwords (
			userid UUID REFERENCES users(id) ON DELETE CASCADE,
			hash BYTEA NOT NULL,
			salt BYTEA NOT NULL,
			algorithm VARCHAR(32) NOT NULL,
			parameters JSONB NOT NULL
		);
		`,
	)

	return err
}

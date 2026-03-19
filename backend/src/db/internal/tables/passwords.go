package tables

import "fmt"

func (q *Tables) InitializePasswordsTable() error {
	// Auth table:
	// id hash salt algorithm parameters
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE passwords (
			userid UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			hash BYTEA NOT NULL,
			salt BYTEA NOT NULL,
			algorithm VARCHAR(32) NOT NULL,
			parameters JSONB NOT NULL
		);
		`,
	)
	if err != nil {
		return fmt.Errorf("could not create passwords table: %v", err)
	}

	return nil
}

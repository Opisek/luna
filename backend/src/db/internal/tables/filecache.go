package tables

import (
	"fmt"
)

func (q *Tables) InitializeFilecacheTable() error {
	var err error
	// Filecache table:
	// id date file
	_, err = q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE filecache (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			date TIMESTAMP NOT NULL,
			name TEXT NOT NULL,
			file BYTEA,
			owner UUID REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create filecache table: %v", err)
	}

	return nil
}

package tables

import (
	"context"
	"fmt"
)

func (q *Tables) InitializeFilecacheTable() error {
	var err error
	// Filecache table:
	// id date file
	_, err = q.Tx.Exec(
		context.TODO(),
		`
		CREATE TABLE IF NOT EXISTS filecache (
			id UUID PRIMARY KEY,
			date TIMESTAMP NOT NULL,
			file BYTEA
		);
	`)
	if err != nil {
		return fmt.Errorf("could not create filecache table: %v", err)
	}

	return nil
}

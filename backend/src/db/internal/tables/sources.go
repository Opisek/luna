package tables

import (
	"context"
	"fmt"
)

func (q *Tables) InitializeSourcesTable() error {
	var err error
	// Sources table:
	// id user name type settings auth
	_, err = q.Tx.Exec(
		context.TODO(),
		`
		CREATE TABLE IF NOT EXISTS sources (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id),
			name VARCHAR(255) NOT NULL,
			type SOURCE_TYPE_ENUM NOT NULL,
			settings JSONB NOT NULL,
			auth_type BYTEA NOT NULL,
			auth BYTEA NOT NULL,
			UNIQUE (user_id, name)
		);
		`,
	)
	if err != nil {
		return fmt.Errorf("could not create sources table: %v", err)
	}

	return nil
}

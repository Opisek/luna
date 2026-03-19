package tables

import "fmt"

func (q *Tables) InitializeInvitesTable() error {
	// Invites table:
	// inviteid author created_at expires_at code
	//
	// UNIQUE on code not only makes sense but also automatically creates an index
	// in postgres
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE invites (
			inviteid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			author UUID REFERENCES users(id) ON DELETE CASCADE,
			email VARCHAR(255) UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at TIMESTAMPTZ NOT NULL,
			code TEXT UNIQUE NOT NULL
		);
		`,
	)
	if err != nil {
		return fmt.Errorf("could not create invites table: %v", err)
	}

	return nil
}

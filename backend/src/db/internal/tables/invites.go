package tables

func (q *Tables) InitializeInvitesTable() error {
	// Invites table:
	// inviteid author created_at expires code
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
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			expires TIMESTAMP NOT NULL,
			code TEXT UNIQUE NOT NULL
		);
		`,
	)

	return err
}

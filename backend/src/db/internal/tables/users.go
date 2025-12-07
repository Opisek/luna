package tables

func (q *Tables) InitializeUsersTable() error {
	// Users table:
	// id username email admin verified enabled searchable pfp created_at
	// TODO: pfp file should reference filecache and be nullable (instead of setting empty UUID)
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,

			admin BOOLEAN,
			verified BOOLEAN,
			enabled BOOLEAN,
			searchable BOOLEAN,

			profile_picture_type PFP_TYPE_ENUM NOT NULL,
			profile_picture_file UUID,
			profile_picture_url VARCHAR(1024),

			created_at TIMESTAMP DEFAULT NOW()
		);
		`,
	)

	return err
}

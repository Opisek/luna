package tables

import "fmt"

func (q *Tables) InitializeTokenPermissionsTable() error {
	// Token permissions table:
	// sessionid permission
	_, err := q.Tx.Exec(
		q.Context,
		`
		CREATE TABLE token_permissions (
			sessionid UUID REFERENCES sessions(sessionid) ON DELETE CASCADE,
			permission VARCHAR(64) NOT NULL
		);
		`,
	)
	if err != nil {
		return fmt.Errorf("could not create tokens permissions table: %v", err)
	}

	_, err = q.Tx.Exec(
		q.Context,
		`
		CREATE INDEX index_token_permissions_sessionid ON token_permissions (sessionid);
	`)
	if err != nil {
		return fmt.Errorf("could not create secondary index on token permissions table: %v", err)
	}

	return nil
}

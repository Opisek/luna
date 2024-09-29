package tables

import (
	"context"
	"fmt"
	"luna-backend/types"
)

type SourceEntry struct {
	Id       types.ID `db:"id" encrypted:"false"`
	UserId   types.ID `db:"userid" encrypted:"false"`
	Name     string   `db:"name" encrypted:"false"`
	Type     string   `db:"type" encrypted:"false"`
	Settings []byte   `db:"settings" encrypted:"false"`
	AuthType string   `db:"auth_type" encrypted:"true"`
	Auth     []byte   `db:"auth" encrypted:"true"`
}

func (q *Tables) InitializeSourcesTable() error {
	var err error
	// Sources table:
	// id user name type settings auth
	_, err = q.Tx.Exec(
		context.TODO(),
		`
		CREATE TABLE IF NOT EXISTS sources (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			userid UUID REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			type SOURCE_TYPE_ENUM NOT NULL,
			settings JSONB NOT NULL,
			auth_type BYTEA NOT NULL,
			auth BYTEA NOT NULL,
			UNIQUE (userid, name)
		);
		`,
	)
	if err != nil {
		return fmt.Errorf("could not create sources table: %v", err)
	}

	return nil
}

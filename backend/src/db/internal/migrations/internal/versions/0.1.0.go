package versions

import (
	"fmt"
	"luna-backend/common"
	"luna-backend/db/internal/migrations/internal/registry"
	"luna-backend/db/internal/migrations/types"
)

func init() {
	registry.RegisterMigration(common.Ver(0, 1, 0), func(q *types.MigrationQueries) error {
		// Support for UUID and encryption
		_, err := q.Tx.Exec(
			q.Context,
			`
			CREATE EXTENSION IF NOT EXISTS pgcrypto;
			`,
		)

		if err != nil {
			return fmt.Errorf("could not create extension pgcrypto: %v", err)
		}

		// Sources enum
		_, err = q.Tx.Exec(
			q.Context,
			`
			CREATE TYPE SOURCE_TYPE_ENUM AS ENUM (
				'caldav',
				'ical'
			);
			`,
		)
		if err != nil {
			return fmt.Errorf("could not create SOURCE_TYPE enum: %v", err)
		}

		// Auth enum
		_, err = q.Tx.Exec(
			q.Context,
			`
			CREATE TYPE AUTH_TYPE_ENUM AS ENUM (
				'none',
				'basic',
				'bearer'
			);
			`,
		)
		if err != nil {
			return fmt.Errorf("could not create AUTH_TYPE enum: %v", err)
		}

		// Tables
		err = q.Tables.InitializeVersionTable()
		if err != nil {
			return fmt.Errorf("could not initialize version table: %v", err)
		}

		err = q.Tables.InitializeUsersTable()
		if err != nil {
			return fmt.Errorf("could not initialize users table: %v", err)
		}

		err = q.Tables.InitializePasswordsTable()
		if err != nil {
			return fmt.Errorf("could not initialize passwords table: %v", err)
		}

		err = q.Tables.InitializeSourcesTable()
		if err != nil {
			return fmt.Errorf("could not initialize sources table: %v", err)
		}

		err = q.Tables.InitializeCalendarsTable()
		if err != nil {
			return fmt.Errorf("could not initialize calendars table: %v", err)
		}

		err = q.Tables.InitializeEventsTable()
		if err != nil {
			return fmt.Errorf("could not initialize events table: %v", err)
		}

		err = q.Tables.InitializeFilecacheTable()
		if err != nil {
			return fmt.Errorf("could not initialize filecache table: %v", err)
		}

		return nil
	})
}

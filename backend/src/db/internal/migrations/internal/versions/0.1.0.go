package versions

import (
	"luna-backend/common"
	"luna-backend/db/internal/migrations/internal/registry"
	"luna-backend/db/internal/migrations/types"
	"luna-backend/errors"
)

func init() {
	registry.RegisterMigration(common.Ver(0, 1, 0), func(q *types.MigrationQueries) *errors.ErrorTrace {
		// Support for UUID and encryption
		_, err := q.Tx.Exec(
			q.Context,
			`
			CREATE EXTENSION IF NOT EXISTS pgcrypto;
			`,
		)

		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not create extension pgcrypto")
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
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not create SOURCE_TYPE enum")
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
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not create AUTH_TYPE enum")
		}

		// Tables
		err = q.Tables.InitializeVersionTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize version table")
		}

		err = q.Tables.InitializeUsersTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize users table")
		}

		err = q.Tables.InitializePasswordsTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize passwords table")
		}

		err = q.Tables.InitializeSourcesTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize sources table")
		}

		err = q.Tables.InitializeCalendarsTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize calendars table")
		}

		err = q.Tables.InitializeEventsTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize events table")
		}

		err = q.Tables.InitializeFilecacheTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize filecache table")
		}

		return nil
	})
}

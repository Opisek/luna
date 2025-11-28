package versions

import (
	"luna-backend/db/internal/migrations/internal/registry"
	migrationTypes "luna-backend/db/internal/migrations/types"
	"luna-backend/errors"
	"luna-backend/types"
)

func init() {
	registry.RegisterMigration(types.Ver(0, 1, 0), func(q *migrationTypes.MigrationQueries) *errors.ErrorTrace {
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

		err = q.Tables.InitializeCalendarOverridesTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize calendar overrides table")
		}

		err = q.Tables.InitializeEventsTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize events table")
		}

		err = q.Tables.InitializeEventOverridesTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize event overrides table")
		}

		err = q.Tables.InitializeFilecacheTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize filecache table")
		}

		err = q.Tables.InitializeUserSettingsTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize user settings table")
		}

		err = q.Tables.InitializeGlobalSettingsTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize global settings table")
		}

		tr := q.Tables.InitializeGlobalSettings()
		if tr != nil {
			return tr
		}

		err = q.Tables.InitializeSessionsTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize sessions table")
		}

		err = q.Tables.InitializeInvitesTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize invites table")
		}

		err = q.Tables.InitializeTokenPermissionsTable()
		if err != nil {
			return errors.New().
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not initialize token permissions table")
		}

		return nil
	})
}

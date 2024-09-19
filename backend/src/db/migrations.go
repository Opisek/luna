package db

import (
	"context"
	"fmt"
	"luna-backend/common"
)

// TODO: add transactions so we revers migrations in case something goes wrong
func (tx *Transaction) RunMigrations(lastVersion *common.Version) error {
	for major := lastVersion.Major; major < len(migrations); major++ {
		for minor := lastVersion.Minor; minor < len(migrations[major]); minor++ {
			for patch := lastVersion.Patch; patch < len(migrations[major][minor]); patch++ {
				migration := migrations[major][minor][patch]
				if migration == nil || major == lastVersion.Major && minor == lastVersion.Minor && patch == lastVersion.Patch {
					continue
				}
				tx.db.logger.Infof("running migration %v.%v.%v", major, minor, patch)
				err := migration(tx)
				if err != nil {
					ver := common.Ver(major, minor, patch)
					return fmt.Errorf("error running migration for %v: %v", ver.String(), err)
				}
			}
		}
	}
	tx.db.logger.Infof("migrations up to date")
	return nil
}

var migrations = [][][]func(*Transaction) error{}

func addMigration(version common.Version, migration func(*Transaction) error) {
	for len(migrations) <= version.Major {
		migrations = append(migrations, [][]func(*Transaction) error{})
	}
	majorMigrations := migrations[version.Major]

	for len(majorMigrations) <= version.Minor {
		majorMigrations = append(majorMigrations, []func(*Transaction) error{})
	}
	migrations[version.Major] = majorMigrations
	minorMigrations := majorMigrations[version.Minor]

	for len(minorMigrations) <= version.Patch {
		minorMigrations = append(minorMigrations, nil)
	}
	migrations[version.Major][version.Minor] = minorMigrations
	migrations[version.Major][version.Minor][version.Patch] = migration
}

func init() {
	// Initialize database
	addMigration(common.Ver(0, 1, 0), func(tx *Transaction) error {
		// Support for UUID and encryption
		_, err := tx.conn.Exec(
			context.TODO(),
			`
			CREATE EXTENSION IF NOT EXISTS pgcrypto;
			`,
		)

		if err != nil {
			return fmt.Errorf("could not create extension pgcrypto: %v", err)
		}

		// Sources enum
		_, err = tx.conn.Exec(
			context.TODO(),
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
		_, err = tx.conn.Exec(
			context.TODO(),
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
		err = tx.initalizeVersionTable()
		if err != nil {
			return fmt.Errorf("could not initialize version table: %v", err)
		}

		err = tx.initializeUserTable()
		if err != nil {
			return fmt.Errorf("could not initialize user table: %v", err)
		}

		err = tx.initializeSourcesTable()
		if err != nil {
			return fmt.Errorf("could not initialize sources table: %v", err)
		}

		return nil
	})

	addMigration(common.Ver(0, 2, 0), func(tx *Transaction) error {
		err := tx.initializeCalendarsTable()
		if err != nil {
			return fmt.Errorf("could not initialize calendars table: %v", err)
		}

		return nil
	})

	addMigration(common.Ver(0, 3, 0), func(tx *Transaction) error {
		err := tx.initializeEventsTable()
		if err != nil {
			return fmt.Errorf("could not initialize events table: %v", err)
		}

		return nil
	})
}

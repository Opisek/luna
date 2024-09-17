package db

import (
	"fmt"
	"luna-backend/common"
)

// TODO: add transactions so we revers migrations in case something goes wrong
func (db *Database) RunMigrations(lastVersion *common.Version) error {
	for major := lastVersion.Major; major < len(migrations); major++ {
		for minor := lastVersion.Minor; minor < len(migrations[major]); minor++ {
			for patch := lastVersion.Patch; patch < len(migrations[major][minor]); patch++ {
				migration := migrations[major][minor][patch]
				if migration == nil || major == lastVersion.Major && minor == lastVersion.Minor && patch == lastVersion.Patch {
					continue
				}
				db.logger.Infof("running migration %v.%v.%v", major, minor, patch)
				err := migration(db)
				if err != nil {
					ver := common.Ver(major, minor, patch)
					return fmt.Errorf("error running migration for %v: %v", ver.String(), err)
				}
			}
		}
	}
	db.logger.Infof("migrations up to date")
	return nil
}

var migrations = [][][]func(*Database) error{}

func addMigration(version common.Version, migration func(*Database) error) {
	for len(migrations) <= version.Major {
		migrations = append(migrations, [][]func(*Database) error{})
	}
	majorMigrations := migrations[version.Major]

	for len(majorMigrations) <= version.Minor {
		majorMigrations = append(majorMigrations, []func(*Database) error{})
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
	addMigration(common.Ver(0, 1, 0), func(db *Database) error {
		// Support for UUID and encryption
		_, err := db.connection.Exec(`
			CREATE EXTENSION IF NOT EXISTS pgcrypto;
		`)

		if err != nil {
			return fmt.Errorf("could not create extension pgcrypto: %v", err)
		}

		// Sources enum
		_, err = db.connection.Exec(`
			CREATE TYPE SOURCE_TYPE_ENUM AS ENUM (
				'caldav',
				'ical'
			);
		`)
		if err != nil {
			return fmt.Errorf("could not create SOURCE_TYPE enum: %v", err)
		}

		// Auth enum
		_, err = db.connection.Exec(`
			CREATE TYPE AUTH_TYPE_ENUM AS ENUM (
				'none',
				'basic',
				'bearer'
			);
		`)
		if err != nil {
			return fmt.Errorf("could not create AUTH_TYPE enum: %v", err)
		}

		// Tables
		err = db.initalizeVersionTable()
		if err != nil {
			return fmt.Errorf("could not initialize version table: %v", err)
		}

		err = db.initializeUserTable()
		if err != nil {
			return fmt.Errorf("could not initialize user table: %v", err)
		}

		err = db.initializeSourcesTable()
		if err != nil {
			return fmt.Errorf("could not initialize sources table: %v", err)
		}

		return nil
	})

	addMigration(common.Ver(0, 2, 0), func(db *Database) error {
		err := db.initializeCalendarsTable()
		if err != nil {
			return fmt.Errorf("could not initialize calendars table: %v", err)
		}

		return nil
	})
}

package db

import (
	"context"
	"fmt"
	"luna-backend/common"
)

func (tx *Transaction) initalizeVersionTable() error {
	// Keeps track of the current backend version as well as stores past
	// versions in case some specific migration rules need to be followed

	// Version table:
	// id major minor patch extension installed

	_, err := tx.conn.Exec(
		context.TODO(),
		`
		CREATE TABLE IF NOT EXISTS version (
			id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			major INT NOT NULL,
			minor INT NOT NULL,
			patch INT NOT NULL,
			extension VARCHAR(255),
			installed TIMESTAMP NOT NULL
		);
		`,
	)

	return err
}

func (tx *Transaction) GetLatestVersion() (common.Version, error) {
	var err error

	err = tx.initalizeVersionTable()
	if err != nil {
		return common.Version{}, fmt.Errorf("could not get latest version: could not initialize version table: %v", err)
	}

	var rowCount int
	err = tx.conn.QueryRow(
		context.TODO(),
		`
		SELECT COUNT(*)
		FROM version;	
		`,
	).Scan(&rowCount)
	if err != nil {
		return common.Version{}, fmt.Errorf("could not get latest version: %v", err)
	}

	if rowCount == 0 {
		return common.EmptyVersion(), nil
	}

	var version common.Version

	err = tx.conn.QueryRow(
		context.TODO(),
		`
		SELECT major, minor, patch, extension
		FROM version
		ORDER BY major DESC, minor DESC, patch DESC
		LIMIT 1
		`,
	).Scan(&version.Major, &version.Minor, &version.Patch, &version.Extension)
	if err != nil {
		return common.EmptyVersion(), fmt.Errorf("could not get latest version: %v", err)
	}

	return version, nil
}

func (tx *Transaction) UpdateVersion(version common.Version) error {
	tx.db.logger.Warnf("updating version to %v", version.String())
	_, err := tx.conn.Exec(
		context.TODO(),
		`
		INSERT INTO version (major, minor, patch, extension, installed)
		VALUES ($1, $2, $3, $4, NOW());
		`,
		version.Major,
		version.Minor,
		version.Patch,
		version.Extension,
	)
	if err != nil {
		return fmt.Errorf("could not update version: %v", err)
	}

	return nil
}

package db

import "luna-backend/common"

func (db *Database) initalizeVersionTable() error {
	// Keeps track of the current backend version as well as stores past
	// versions in case some specific migration rules need to be followed

	// Version table:
	// id major minor patch extension installed

	_, err := db.connection.Exec(`
		CREATE TABLE IF NOT EXISTS version (
			id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			major INT NOT NULL,
			minor INT NOT NULL,
			patch INT NOT NULL,
			extension VARCHAR(255),
			installed TIMESTAMP NOT NULL
		);
	`)

	return err
}

func (db *Database) GetLatestVersion() (common.Version, error) {
	var err error

	err = db.initalizeVersionTable()
	if err != nil {
		db.logger.Errorf("could not initialize version table: %v", err)
		return common.Version{}, err
	}

	var rowCount int
	err = db.connection.QueryRow(`
		SELECT COUNT(*)
		FROM version;	
	`).Scan(&rowCount)
	if err != nil {
		db.logger.Errorf("could not get latest version: %v", err)
		return common.Version{}, err
	}

	if rowCount == 0 {
		return common.EmptyVersion(), nil
	}

	var version common.Version

	err = db.connection.QueryRow(`
		SELECT major, minor, patch, extension
		FROM version
		ORDER BY major DESC, minor DESC, patch DESC
		LIMIT 1
	`).Scan(&version.Major, &version.Minor, &version.Patch, &version.Extension)
	if err != nil {
		db.logger.Errorf("could not get latest version: %v", err)
		return common.EmptyVersion(), err
	}

	return version, nil
}

func (db *Database) UpdateVersion(version common.Version) error {
	db.logger.Warnf("updating version to %v", version.String())
	_, err := db.connection.Exec(`
		INSERT INTO version (major, minor, patch, extension, installed)
		VALUES ($1, $2, $3, $4, NOW());
	`, version.Major, version.Minor, version.Patch, version.Extension)
	if err != nil {
		db.logger.Errorf("could not update version: %v", err)
		return err
	}

	return nil
}

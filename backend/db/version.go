package db

import "luna-backend/common"

func (db *Database) GetLatestVersion() (common.Version, error) {
	var err error

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
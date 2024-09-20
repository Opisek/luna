package queries

import (
	"context"
	"fmt"
	"luna-backend/common"
)

func (q *Queries) GetLatestVersion() (common.Version, error) {
	var err error

	var rowCount int
	err = q.Tx.QueryRow(
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

	err = q.Tx.QueryRow(
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

func (q *Queries) UpdateVersion(version common.Version) error {
	q.Logger.Warnf("updating version to %v", version.String())
	_, err := q.Tx.Exec(
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

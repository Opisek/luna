package queries

import (
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"
)

func (q *Queries) GetLatestVersion() (types.Version, *errors.ErrorTrace) {
	var err error

	var rowCount int
	err = q.Tx.QueryRow(
		q.Context,
		`
		SELECT COUNT(*)
		FROM version;	
		`,
	).Scan(&rowCount)
	if err != nil {
		return types.Version{}, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not check saved versions").
			Append(errors.LvlPlain, "Database error")
	}

	if rowCount == 0 {
		return types.EmptyVersion(), nil
	}

	var version types.Version

	err = q.Tx.QueryRow(
		q.Context,
		`
		SELECT major, minor, patch, extension
		FROM version
		ORDER BY major DESC, minor DESC, patch DESC
		LIMIT 1
		`,
	).Scan(&version.Major, &version.Minor, &version.Patch, &version.Extension)
	if err != nil {
		return types.EmptyVersion(), errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not get last used version").
			Append(errors.LvlPlain, "Database error")
	}

	return version, nil
}

func (q *Queries) UpdateVersion(version types.Version) *errors.ErrorTrace {
	q.Logger.Warnf("updating version to %v", version.String())
	_, err := q.Tx.Exec(
		q.Context,
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
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not update used version").
			Append(errors.LvlPlain, "Database error")
	}

	return nil
}

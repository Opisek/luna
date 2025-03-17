package queries

import (
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// Password errors are kept vague for security purposes.

func (q *Queries) GetPassword(userId types.ID) (*types.PasswordEntry, *errors.ErrorTrace) {
	var err error

	entry := &types.PasswordEntry{}

	err = q.Tx.QueryRow(
		q.Context,
		`
		SELECT hash, salt, algorithm, parameters
		FROM passwords
		WHERE userid = $1;
		`, userId.UUID(),
	).Scan(&entry.Hash, &entry.Salt, &entry.Algorithm, &entry.Parameters)

	switch err {
	case nil:
		return entry, nil
	case pgx.ErrNoRows:
		fallthrough // for security reasons, we should not return overly detailed errors here
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get password information for user %v", userId)
	}
}

func (q *Queries) InsertPassword(userId types.ID, entry *types.PasswordEntry) *errors.ErrorTrace {
	var err error

	_, err = q.Tx.Exec(
		q.Context,
		`
		INSERT INTO passwords (userid, hash, salt, algorithm, parameters)
		VALUES ($1, $2, $3, $4, $5);
		`,
		userId.UUID(), entry.Hash, entry.Salt, entry.Algorithm, entry.Parameters,
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not insert password information for user %v", userId)
	}

	return nil
}

func (q *Queries) UpdatePassword(userId types.ID, entry *types.PasswordEntry) *errors.ErrorTrace {
	var err error

	_, err = q.Tx.Exec(
		q.Context,
		`
		UPDATE passwords
		SET hash = $1, salt = $2, algorithm = $3, parameters = $4
		WHERE userid = $5;
		`,
		entry.Hash, entry.Salt, entry.Algorithm, entry.Parameters,
		userId.UUID(),
	)

	switch err {
	case nil:
		return nil
	case pgx.ErrNoRows:
		fallthrough // for security reasons, we should not return overly detailed errors here
	default:
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not update password information for user %v", userId)
	}
}

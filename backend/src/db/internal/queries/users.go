package queries

import (
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/google/uuid"
)

// As with password-related queries, user management errors are kept vague on purpose.

func (q *Queries) AddUser(user *types.User) (types.ID, *errors.ErrorTrace) {
	var err error

	query := `
		INSERT INTO users (username, email, admin, searchable, profile_picture)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var id types.ID
	err = q.Tx.QueryRow(q.Context, query, user.Username, user.Email, user.Admin, user.Searchable, user.ProfilePicture).Scan(&id)

	if err != nil {
		return types.EmptyId(), errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not add user")
	}

	user.Id = id

	return id, nil
}

func (q *Queries) GetUserIdFromEmail(email string) (types.ID, *errors.ErrorTrace) {
	var err error

	var id uuid.UUID

	err = q.Tx.QueryRow(
		q.Context,
		`
		SELECT id
		FROM users
		WHERE email = $1;
		`,
		email,
	).Scan(&id)

	if err != nil {
		return types.EmptyId(), errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get user with email %v", email)
	}

	return types.IdFromUuid(id), nil
}

func (q *Queries) GetUserIdFromUsername(username string) (types.ID, *errors.ErrorTrace) {
	var err error

	var id uuid.UUID

	err = q.Tx.QueryRow(
		q.Context,
		`
		SELECT id
		FROM users
		WHERE username = $1;
		`,
		username,
	).Scan(&id)

	if err != nil {
		return types.EmptyId(), errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get user with username %v", username)
	}
	return types.IdFromUuid(id), nil
}

func (q *Queries) IsAdmin(userId types.ID) (bool, *errors.ErrorTrace) {
	var err error

	var admin bool

	err = q.Tx.QueryRow(
		q.Context,
		`
		SELECT admin
		FROM users
		WHERE id = $1;
		`,
		userId.UUID(),
	).Scan(&admin)

	if err != nil {
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not check admin status of user %v", userId)
	}

	return admin, nil
}

func (q *Queries) AnyUsersExist() (bool, *errors.ErrorTrace) {
	rows, err := q.Tx.Query(
		q.Context,
		`
		SELECT *
		FROM users
		LIMIT 1;
		`,
	)

	if err != nil {
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not check the presence of any users")
	}

	exists := rows.Next()
	rows.Close()

	return exists, nil
}

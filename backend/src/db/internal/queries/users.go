package queries

import (
	"luna-backend/constants"
	"luna-backend/errors"
	"luna-backend/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// As with password-related queries, user management errors are kept vague on purpose.

func (q *Queries) AddUser(user *types.User) (types.ID, *errors.ErrorTrace) {
	var err error

	query := `
		INSERT INTO users (username, email, admin, verified, enabled, searchable, profile_picture_type, profile_picture_file, profile_picture_url)
		VALUES ($1, $2, $3, FALSE, TRUE, $4, $5, $6, $7)
		RETURNING id;
	`

	var id types.ID
	err = q.Tx.QueryRow(q.Context, query, user.Username, user.Email, user.Admin, user.Searchable, user.ProfilePictureType, user.ProfilePictureFile, user.ProfilePictureUrl).Scan(&id)

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

	switch err {
	case nil:
		return admin, nil
	case pgx.ErrNoRows:
		return false, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "User %v does not exist", userId).
			AltStr(errors.LvlPlain, "User does not exist")
	default:
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not check admin status of user %v", userId)
	}
}

func (q *Queries) AnyUsersExist() (bool, *errors.ErrorTrace) {
	// TODO: rewrite with EXISTS
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

func (q *Queries) GetUser(userId types.ID) (*types.User, *errors.ErrorTrace) {
	var err error

	user := &types.User{}
	var rawProfilePictureUrl string

	err = q.Tx.QueryRow(
		q.Context,
		`
		SELECT id, username, email, admin, verified, enabled, searchable, profile_picture_type, COALESCE(profile_picture_file, uuid_nil()), COALESCE(profile_picture_url, ''), created_at
		FROM users
		WHERE id = $1;
		`,
		userId.UUID(),
	).Scan(&user.Id, &user.Username, &user.Email, &user.Admin, &user.Verified, &user.Enabled, &user.Searchable, &user.ProfilePictureType, &user.ProfilePictureFile, &rawProfilePictureUrl, &user.CreatedAt)
	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return nil, errors.New().Status(http.StatusNotFound).
			Append(errors.LvlDebug, "User %v does not exist", userId).
			Append(errors.LvlDebug, "Could not get user %v", userId).
			AltStr(errors.LvlPlain, "User not found")
	default:
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Database encountered an error").
			Append(errors.LvlDebug, "Could not get user %v", userId).
			AltStr(errors.LvlPlain, "Database error")
	}

	if rawProfilePictureUrl == "" && user.ProfilePictureType != constants.ProfilePictureDatabase {
		return nil, errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "The profile picture URL is empty").
			Append(errors.LvlDebug, "Could not get user %v", userId).
			AltStr(errors.LvlPlain, "Database error")
	}
	user.ProfilePictureUrl, err = types.NewUrl(rawProfilePictureUrl)
	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not parse the profile picture URL").
			Append(errors.LvlDebug, "Could not get user %v", userId).
			AltStr(errors.LvlPlain, "Database error")
	}

	return user, nil
}

func (q *Queries) GetUsers(all bool) ([]*types.User, *errors.ErrorTrace) {
	var err error

	var rows pgx.Rows

	var query string
	if all {
		query = `
		SELECT id, username, email, admin, verified, enabled, searchable, profile_picture_type, COALESCE(profile_picture_file, uuid_nil()), COALESCE(profile_picture_url, ''), created_at
		FROM users;
		`
	} else {
		query = `
		SELECT id, username, email, admin, verified, enabled, searchable, profile_picture_type, COALESCE(profile_picture_file, uuid_nil()), COALESCE(profile_picture_url, ''), created_at
		FROM users
		WHERE enabled = TRUE
		AND searchable = TRUE;
		`
	}

	rows, err = q.Tx.Query(
		q.Context,
		query,
	)

	if err != nil {
		return nil, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not get users")
	}

	defer rows.Close()

	users := make([]*types.User, 0)

	for rows.Next() {
		user := &types.User{}
		var rawProfilePictureUrl string

		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Admin, &user.Verified, &user.Enabled, &user.Searchable, &user.ProfilePictureType, &user.ProfilePictureFile, &rawProfilePictureUrl, &user.CreatedAt)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not get users").
				Append(errors.LvlPlain, "Database error")
		}

		user.ProfilePictureUrl, err = types.NewUrl(rawProfilePictureUrl)
		if err != nil {
			return nil, errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not get user %v", user.Id).
				Append(errors.LvlPlain, "Database error")
		}

		users = append(users, user)
	}

	return users, nil
}

func (q *Queries) UpdateUserData(user *types.User) *errors.ErrorTrace {
	var err error

	query := `
		UPDATE users
		SET username = $1, email = $2, searchable = $3, profile_picture_type = $4, profile_picture_file = $5, profile_picture_url = $6
		WHERE id = $7;
	`

	_, err = q.Tx.Exec(
		q.Context,
		query,
		user.Username,
		user.Email,
		user.Searchable,
		user.ProfilePictureType,
		user.ProfilePictureFile,
		user.ProfilePictureUrl,
		user.Id.UUID(),
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not update user %v", user.Id).
			AltStr(errors.LvlPlain, "Database error")
	}

	return nil
}

func (q *Queries) DeleteUser(userId types.ID) *errors.ErrorTrace {
	var err error

	query := `
		DELETE FROM users
		WHERE id = $1;
	`

	_, err = q.Tx.Exec(
		q.Context,
		query,
		userId.UUID(),
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not delete user %v", userId).
			AltStr(errors.LvlPlain, "Database error")
	}

	return nil
}

func (q *Queries) IsUserEnabled(userId types.ID) (bool, *errors.ErrorTrace) {
	var err error

	var enabled bool

	err = q.Tx.QueryRow(
		q.Context,
		`
		SELECT enabled
		FROM users
		WHERE id = $1;
		`,
		userId.UUID(),
	).Scan(&enabled)

	if err != nil {
		return false, errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not check if user %v is enabled", userId)
	}

	return enabled, nil
}

func (q *Queries) SetUserEnabled(userId types.ID, enabled bool) *errors.ErrorTrace {
	var err error

	query := `
		UPDATE users
		SET enabled = $1
		WHERE id = $2
		AND admin = FALSE;
	`

	_, err = q.Tx.Exec(
		q.Context,
		query,
		enabled,
		userId.UUID(),
	)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not set user %v enabled to %v", userId, enabled)
	}

	return nil
}

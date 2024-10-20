package queries

import (
	"context"
	"fmt"
	"luna-backend/types"

	"github.com/google/uuid"
)

func (q *Queries) AddUser(user *types.User) (types.ID, error) {
	var err error

	query := `
		INSERT INTO users (username, email, admin)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var id types.ID
	err = q.Tx.QueryRow(context.TODO(), query, user.Username, user.Email, user.Admin).Scan(&id)
	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not add user: %v", err)
	}

	user.Id = id

	return id, err
}

func (q *Queries) GetUserIdFromEmail(email string) (uuid.UUID, error) {
	var err error

	var id uuid.UUID

	err = q.Tx.QueryRow(
		context.TODO(),
		`
		SELECT id
		FROM users
		WHERE email = $1;
		`,
		email,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get user id by email %v: %v", email, err)
	}
	return id, err
}

func (q *Queries) GetUserIdFromUsername(username string) (types.ID, error) {
	var err error

	var id uuid.UUID

	err = q.Tx.QueryRow(
		context.TODO(),
		`
		SELECT id
		FROM users
		WHERE username = $1;
		`,
		username,
	).Scan(&id)
	if err != nil {
		return types.EmptyId(), fmt.Errorf("could not get user id by username %v: %v", username, err)
	}
	return types.IdFromUuid(id), err
}

func (q *Queries) IsAdmin(id int) (bool, error) {
	var err error

	var admin bool

	err = q.Tx.QueryRow(
		context.TODO(),
		`
		SELECT admin
		FROM users
		WHERE id = $1;
		`,
		id,
	).Scan(&admin)

	if err != nil {
		return false, fmt.Errorf("could not get admin status of user %v: %v", id, err)
	}
	return admin, err
}

func (q *Queries) AnyUsersExist() (bool, error) {
	rows, err := q.Tx.Query(
		context.TODO(),
		`
		SELECT *
		FROM users
		LIMIT 1;
		`,
	)

	if err != nil {
		return false, fmt.Errorf("could not check if any users exist: %v", err)
	}

	exists := rows.Next()
	rows.Close()

	return exists, nil
}

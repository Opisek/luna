package queries

import (
	"context"
	"fmt"
	"luna-backend/types"

	"github.com/google/uuid"
)

// TODO: return the created user's ID
func (q *Queries) AddUser(user *types.User) error {
	var err error

	_, err = q.Tx.Exec(
		context.TODO(),
		`
		INSERT INTO users (username, password, algorithm, email, admin)
		VALUES ($1, $2, $3, $4, $5);
		`,
		user.Username,
		user.Password,
		user.Algorithm,
		user.Email,
		user.Admin,
	)
	if err != nil {
		return fmt.Errorf("could not add user: %v", err)
	}
	return nil
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

func (q *Queries) GetPassword(id types.ID) (string, string, error) {
	var err error

	var password, algorithm string

	err = q.Tx.QueryRow(
		context.TODO(),
		`
		SELECT password, algorithm
		FROM users
		WHERE id = $1;
		`, id.UUID(),
	).Scan(&password, &algorithm)

	if err != nil {
		return "", "", fmt.Errorf("could not get password hash of user %v: %v", id, err)
	}
	return password, algorithm, err
}

func (q *Queries) UpdatePassword(id types.ID, password string, alg string) error {
	var err error

	_, err = q.Tx.Exec(
		context.TODO(),
		`
		UPDATE users
		SET password = $1, algorithm = $2
		WHERE id = $3;
		`,
		password,
		alg,
		id.UUID(),
	)

	if err != nil {
		return fmt.Errorf("could not update password of user %v: %v", id, err)
	}
	return err
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

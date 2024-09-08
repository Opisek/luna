package db

import (
	"fmt"
	"luna-backend/types"

	"github.com/google/uuid"
)

func (db *Database) initializeUserTable() error {
	// Auth table:
	// id username password email admin
	_, err := db.connection.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			algorithm VARCHAR(32) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			admin BOOLEAN
		);
	`)

	return err
}

// TODO: return the created user's ID
func (db *Database) AddUser(user *types.User) error {
	var err error

	_, err = db.connection.Exec(`
		INSERT INTO users (username, password, algorithm, email, admin)
		VALUES ($1, $2, $3, $4, $5);
	`, user.Username, user.Password, user.Algorithm, user.Email, user.Admin)
	if err != nil {
		return fmt.Errorf("could not add user: %v", err)
	}
	return nil
}

func (db *Database) GetUserIdFromEmail(email string) (uuid.UUID, error) {
	var err error

	var id uuid.UUID

	err = db.connection.QueryRow(`
		SELECT id
		FROM users
		WHERE email = $1;
	`, email).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get user id by email %v: %v", email, err)
	}
	return id, err
}

func (db *Database) GetUserIdFromUsername(username string) (uuid.UUID, error) {
	var err error

	var id uuid.UUID

	err = db.connection.QueryRow(`
		SELECT id
		FROM users
		WHERE username = $1;
	`, username).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not get user id by username %v: %v", username, err)
	}
	return id, err
}

func (db *Database) GetPassword(id uuid.UUID) (string, string, error) {
	var err error

	var password, algorithm string

	err = db.connection.QueryRow(`
		SELECT password, algorithm
		FROM users
		WHERE id = $1;
	`, id).Scan(&password, &algorithm)

	if err != nil {
		return "", "", fmt.Errorf("could not get password hash of user %v: %v", id, err)
	}
	return password, algorithm, err
}

func (db *Database) UpdatePassword(id uuid.UUID, password string, alg string) error {
	var err error

	_, err = db.connection.Exec(`
		UPDATE users
		SET password = $1, algorithm = $2
		WHERE id = $3;
	`, password, alg, id)

	if err != nil {
		return fmt.Errorf("could not update password of user %v: %v", id, err)
	}
	return err
}

func (db *Database) IsAdmin(id int) (bool, error) {
	var err error

	var admin bool

	err = db.connection.QueryRow(`
		SELECT admin
		FROM users
		WHERE id = $1;
	`, id).Scan(&admin)

	if err != nil {
		return false, fmt.Errorf("could not get admin status of user %v: %v", id, err)
	}
	return admin, err
}

func (db *Database) AnyUsersExist() (bool, error) {
	rows, err := db.connection.Query(`
		SELECT *
		FROM users
		LIMIT 1;
	`)

	if err != nil {
		return false, fmt.Errorf("could not check if any users exist: %v", err)
	}

	exists := rows.Next()
	rows.Close()

	return exists, nil
}

package db

import "luna-backend/types"

// TODO: return the created user's ID
func (db *Database) AddUser(user *types.User) error {
	var err error

	_, err = db.connection.Exec(`
		INSERT INTO users (username, password, algorithm, email, admin)
		VALUES ($1, $2, $3, $4, $5);
	`, user.Username, user.Password, user.Algorithm, user.Email, user.Admin)
	if err != nil {
		db.logger.Errorf("could not add user: %v", err)
	}
	return err
}

func (db *Database) GetUserIdFromEmail(email string) (int, error) {
	var err error

	var id int

	err = db.connection.QueryRow(`
		SELECT id
		FROM users
		WHERE email = $1;
	`, email).Scan(&id)
	if err != nil {
		db.logger.Errorf("could not get user id by email %v: %v", email, err)
	}
	return id, err
}

func (db *Database) GetUserIdFromUsername(username string) (int, error) {
	var err error

	var id int

	err = db.connection.QueryRow(`
		SELECT id
		FROM users
		WHERE username = $1;
	`, username).Scan(&id)
	if err != nil {
		db.logger.Errorf("could not get user id by username %v: %v", username, err)
	}
	return id, err
}

func (db *Database) GetPassword(id int) (string, string, error) {
	var err error

	var password, algorithm string

	err = db.connection.QueryRow(`
		SELECT password, algorithm
		FROM users
		WHERE id = $1;
	`, id).Scan(&password, &algorithm)

	if err != nil {
		db.logger.Errorf("could not get password hash of user %v: %v", id, err)
	}
	return password, algorithm, err
}

func (db *Database) UpdatePassword(id int, password string, alg string) error {
	var err error

	_, err = db.connection.Exec(`
		UPDATE users
		SET password = $1, algorithm = $2
		WHERE id = $3;
	`, password, alg, id)

	if err != nil {
		db.logger.Errorf("could not update password of user %v: %v", id, err)
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
		db.logger.Errorf("could not get admin status of user %v: %v", id, err)
	}
	return admin, err
}

func (db *Database) AnyUsersExist() bool {
	row := db.connection.QueryRow(`
		SELECT *
		FROM users
		LIMIT 1;
	`)

	return row != nil
}

package db

// This only initializes tables, it does not handle migrations
func (db *Database) InitializeTables() error {
	var err error

	err = db.initalizeVersionTable()
	if err != nil {
		db.logger.Errorf("could not initialize version table: %v", err)
		return err
	}

	err = db.initializeUserTable()
	if err != nil {
		db.logger.Errorf("could not initialize user table: %v", err)
		return err
	}

	return nil
}

func (db *Database) initalizeVersionTable() error {
	// Keeps track of the current backend version as well as stores past
	// versions in case some specific migration rules need to be followed

	// Version table:
	// id major minor patch extension installed

	_, err := db.connection.Exec(`
		CREATE TABLE IF NOT EXISTS version (
			id SERIAL PRIMARY KEY,
			major INT NOT NULL,
			minor INT NOT NULL,
			patch INT NOT NULL,
			extension VARCHAR(255),
			installed TIMESTAMP NOT NULL
		);
	`)

	return err
}

func (db *Database) initializeUserTable() error {
	// Auth table:
	// id username password email admin
	_, err := db.connection.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			algorithm VARCHAR(32) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			admin BOOLEAN
		);
	`)

	return err
}

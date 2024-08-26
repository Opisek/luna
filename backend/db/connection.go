package db

import (
	"github.com/jackc/pgx"
)

func (db *Database) Connect() error {
	if db.connection != nil && db.connection.IsAlive() {
		return nil
	}

	var err error
	// TODO: add connection pool
	db.connection, err = pgx.Connect(*db.pgxConfig)

	if err != nil {
		db.logger.Errorf("could not connect to database: %v", err)
	}

	return err
}

func (db *Database) Disconnect() {
	db.connection.Close()
}

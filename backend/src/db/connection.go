package db

import (
	"errors"
	"time"

	"github.com/jackc/pgx"
)

var retries int = 5

func (db *Database) IsConnected() bool {
	return db.connection != nil && db.connection.IsAlive()
}

func (db *Database) Connect() error {
	if db.IsConnected() {
		return nil
	}

	for i := 0; i < retries; i++ {
		err := db.connect()

		if err == nil {
			return nil
		}

		if i < retries-1 {
			db.logger.Warnf("could not connect to database: %v\nretrying...", err)
			time.Sleep(1 * time.Second)
		} else {
			db.logger.Errorf("could not connect to database: %v", err)
		}
	}

	return errors.New("could not connect to database")
}

func (db *Database) connect() error {
	var err error
	db.connection, err = pgx.Connect(*db.pgxConfig)
	return err
}

func (db *Database) Disconnect() {
	db.connection.Close()
}

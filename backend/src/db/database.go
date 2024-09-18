package db

import (
	"context"
	"fmt"

	"luna-backend/common"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var retries int = 5

type Database struct {
	pgxConfig *pgx.ConnConfig
	pool      *pgxpool.Pool

	commonConfig *common.CommonConfig
	logger       *logrus.Entry
}

func NewDatabase(host string, port uint16, username, password, database string, commonConfig *common.CommonConfig, logger *logrus.Entry) *Database {
	pgxConfig := &pgx.ConnConfig{
		Host:     host,
		Port:     port,
		User:     username,
		Password: password,
		Database: database,
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, port, database)

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		logger.Fatalf("could not create database pool: %v", err)
	}

	db := &Database{
		pgxConfig:    pgxConfig,
		pool:         pool,
		commonConfig: commonConfig,
		logger:       logger,
	}

	return db
}

//func (db *Database) IsConnected() bool {
//	return tx.conn != nil && db.connection.IsAlive()
//}

//func (db *Database) Connect() error {
//	if db.IsConnected() {
//		return nil
//	}
//
//	for i := 0; i < retries; i++ {
//		err := db.connect()
//
//		if err == nil {
//			return nil
//		}
//
//		if i < retries-1 {
//			db.logger.Warnf("could not connect to database: %v\nretrying...", err)
//			time.Sleep(1 * time.Second)
//		} else {
//			db.logger.Errorf("could not connect to database: %v", err)
//		}
//	}
//
//	return errors.New("could not connect to database")
//}
//
//func (db *Database) connect() error {
//	var err error
//	db.connection, err = pgx.Connect(*db.pgxConfig)
//	return err
//}
//
//func (db *Database) Disconnect() {
//	db.connection.Close()
//}

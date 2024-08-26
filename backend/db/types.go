package db

import (
	"luna-backend/common"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type Database struct {
	pgxConfig  *pgx.ConnConfig
	connection *pgx.Conn

	commonConfig *common.CommonConfig
	logger       *logrus.Entry
}

func NewDatabase(host string, port uint16, username, password, database string, commonConfig *common.CommonConfig, logger *logrus.Entry) *Database {
	db := &Database{
		pgxConfig: &pgx.ConnConfig{
			Host:     host,
			Port:     port,
			User:     username,
			Password: password,
			Database: database,
		},
		commonConfig: commonConfig,
		logger:       logger,
	}

	return db
}

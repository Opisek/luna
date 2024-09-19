package db

import (
	"context"
	"fmt"

	"luna-backend/common"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

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

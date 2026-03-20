package db

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"luna-backend/config"
	"luna-backend/db/internal/parsing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Database struct {
	pool *pgxpool.Pool

	commonConfig     *config.CommonConfig
	primitivesParser parsing.PrimitivesParser

	logger *logrus.Entry
}

func NewDatabase(connStr string, host string, port uint16, username, password, database string, commonConfig *config.CommonConfig, primitivesParser parsing.PrimitivesParser, logger *logrus.Entry) *Database {

	if connStr == "" {
		connStr = (&url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(username, password),
			Host:   net.JoinHostPort(host, fmt.Sprintf("%d", port)),
			Path:   database,
		}).String()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Fatalf("could not create database pool: %v", err)
	}

	db := &Database{
		pool:             pool,
		commonConfig:     commonConfig,
		primitivesParser: primitivesParser,
		logger:           logger,
	}

	return db
}

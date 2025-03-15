package db

import (
	"context"
	"fmt"
	"time"

	"luna-backend/common"
	"luna-backend/db/internal/parsing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Database struct {
	pool *pgxpool.Pool

	commonConfig     *common.CommonConfig
	primitivesParser parsing.PrimitivesParser

	logger *logrus.Entry
}

func NewDatabase(host string, port uint16, username, password, database string, commonConfig *common.CommonConfig, primitivesParser parsing.PrimitivesParser, logger *logrus.Entry) *Database {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, port, database)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, url)
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

package db

import (
	"context"
	"fmt"
	"luna-backend/db/internal/migrations"
	"luna-backend/db/internal/migrations/types"
	"luna-backend/db/internal/queries"
	"luna-backend/db/internal/tables"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Transaction struct {
	db *Database
	tx pgx.Tx

	queries    *queries.Queries
	tables     *tables.Tables
	migrations *types.MigrationQueries
}

func (db *Database) BeginTransaction() (*Transaction, error) {
	tx, err := db.pool.Begin(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}

	transaction := &Transaction{
		db: db,
		tx: tx,
	}

	return transaction, nil
}

func (tx *Transaction) Commit(logger *logrus.Entry) error {
	err := tx.tx.Commit(context.TODO())

	if err != nil {
		err := fmt.Errorf("could not commit transaction: %v", err)
		logger.Error(err)
		return err
	}

	return nil
}

func (tx *Transaction) Rollback(logger *logrus.Entry) error {
	err := tx.tx.Rollback(context.TODO())

	if err != nil && err != pgx.ErrTxClosed {
		err := fmt.Errorf("could not rollback transaction: %v", err)
		logger.Error(err)
		return err
	}

	return nil
}

func (tx *Transaction) Queries() *queries.Queries {
	if tx.queries == nil {
		tx.queries = &queries.Queries{
			Tx:           tx.tx,
			Logger:       tx.db.logger,
			CommonConfig: tx.db.commonConfig,
		}
	}
	return tx.queries
}

func (tx *Transaction) Tables() *tables.Tables {
	if tx.tables == nil {
		tx.tables = &tables.Tables{
			Tx: tx.tx,
		}
	}
	return tx.tables
}

func (tx *Transaction) Migrations() *types.MigrationQueries {
	if tx.migrations == nil {
		tx.migrations = migrations.NewMigrationQueries(tx.tx, tx.db.logger, tx.db.commonConfig, tx.Tables())
	}
	return tx.migrations
}

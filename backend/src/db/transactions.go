package db

import (
	"context"
	"fmt"
	"luna-backend/db/internal/migrations"
	"luna-backend/db/internal/migrations/types"
	"luna-backend/db/internal/queries"
	"luna-backend/db/internal/tables"
	"luna-backend/errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Transaction struct {
	db      *Database
	context context.Context
	tx      pgx.Tx

	queries    *queries.Queries
	tables     *tables.Tables
	migrations *types.MigrationQueries
}

func (db *Database) BeginTransaction(ctx context.Context) (*Transaction, *errors.ErrorTrace) {
	tx, err := db.pool.Begin(ctx)

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	if err != nil {
		fmt.Println(errMsg)
	}
	switch {
	case err == nil:
		break
	case strings.Contains(errMsg, "connection refused"):
		return nil, errors.New().Status(http.StatusServiceUnavailable).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlPlain, "The database is not reachable")
	case strings.Contains(errMsg, "authentication failed"):
		return nil, errors.New().Status(http.StatusServiceUnavailable).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Wrong database credentials").
			AltStr(errors.LvlPlain, "Database error")
	default:
		return nil, errors.New().Status(http.StatusServiceUnavailable).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not begin transaction").
			AltStr(errors.LvlPlain, "Database error")
	}

	transaction := &Transaction{
		db:      db,
		context: ctx,
		tx:      tx,
	}

	return transaction, nil
}

func (tx *Transaction) Commit(logger *logrus.Entry) *errors.ErrorTrace {
	err := tx.tx.Commit(tx.context)

	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not commit transaction").
			AltStr(errors.LvlPlain, "Database error")
	}

	return nil
}

func (tx *Transaction) Rollback(logger *logrus.Entry) *errors.ErrorTrace {
	err := tx.tx.Rollback(tx.context)

	if err != nil && err != pgx.ErrTxClosed && !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlWordy, "Could not rollback transaction").
			AltStr(errors.LvlPlain, "Database error")
	}

	return nil
}

func (tx *Transaction) Queries() *queries.Queries {
	if tx.queries == nil {
		tx.queries = &queries.Queries{
			Tx:               tx.tx,
			Context:          tx.context,
			Logger:           tx.db.logger,
			CommonConfig:     tx.db.commonConfig,
			PrimitivesParser: &tx.db.primitivesParser,
		}
	}
	return tx.queries
}

func (tx *Transaction) Tables() *tables.Tables {
	if tx.tables == nil {
		tx.tables = &tables.Tables{
			Tx:      tx.tx,
			Context: tx.context,
		}
	}
	return tx.tables
}

func (tx *Transaction) Migrations() *types.MigrationQueries {
	if tx.migrations == nil {
		tx.migrations = migrations.NewMigrationQueries(tx.tx, tx.context, tx.db.logger, tx.db.commonConfig, tx.Tables())
	}
	return tx.migrations
}

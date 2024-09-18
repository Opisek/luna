package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Transaction struct {
	db   *Database
	conn pgx.Tx
}

func (db *Database) BeginTransaction() (*Transaction, error) {
	tx, err := db.pool.Begin(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %v", err)
	}

	transaction := &Transaction{
		db:   db,
		conn: tx,
	}

	return transaction, nil
}

func (tx *Transaction) Commit(logger *logrus.Entry) error {
	err := tx.conn.Commit(context.TODO())

	if err != nil {
		err := fmt.Errorf("could not commit transaction: %v", err)
		logger.Error(err)
		return err
	}

	return nil
}

func (tx *Transaction) Rollback(logger *logrus.Entry) error {
	err := tx.conn.Rollback(context.TODO())

	if err != nil && err != pgx.ErrTxClosed {
		err := fmt.Errorf("could not rollback transaction: %v", err)
		logger.Error(err)
		return err
	}

	return nil
}

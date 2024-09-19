package db

import (
	"context"
	"fmt"
	"luna-backend/crypto"

	"github.com/jackc/pgx/v5"
)

func (tx *Transaction) CopyAndUpdate(context context.Context, tableName string, columnNames []string, rows [][]any) error {
	randomNumber, err := crypto.GenerateRandomNumber()
	if err != nil {
		return fmt.Errorf("could not copy into table %v: %v", tableName, err)
	}

	tpmTableName := fmt.Sprintf("temp_%v_%v", tableName, randomNumber)

	tx.conn.Exec(
		context,
		`
		CREATE TEMPORARY TABLE $1 (
			LIKE $2 INCLUDING ALL	
		)	ON COMMIT DELET ROWS;
		`,
		tpmTableName,
		tableName,
	)

	_, err = tx.conn.CopyFrom(
		context,
		pgx.Identifier{tpmTableName},
		columnNames,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("could not copy into table %v: %v", tableName, err)
	}

	_, err = tx.conn.Exec(
		context,
		`
		INSERT INTO $1
		SELECT *
		FROM $2
		WITH ON CONFLICT DO UPDATE;
		`,
		tableName,
		tpmTableName,
	)
	if err != nil {
		return fmt.Errorf("could not copy into table %v: %v", tableName, err)
	}

	return nil
}

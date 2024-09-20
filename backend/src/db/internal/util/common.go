package util

import (
	"context"
	"errors"
	"fmt"
	"luna-backend/crypto"
	"strings"

	"github.com/jackc/pgx/v5"
)

func CopyAndUpdate(tx pgx.Tx, context context.Context, tableName string, columnNames []string, updateColumns []string, rows [][]any) error {
	randomNumber, err := crypto.GenerateRandomNumber()
	if err != nil {
		return fmt.Errorf("could not copy into table %v: %v", tableName, err)
	}

	if !isSafe(tableName) {
		return errors.New("could not copy into table: table name contains illegal substrings")
	}
	tpmTableName := fmt.Sprintf("temp_%v_%v", tableName, randomNumber)

	query := fmt.Sprintf(
		`
		CREATE TEMP TABLE %s (
			LIKE %s INCLUDING ALL	
		)	ON COMMIT DELETE ROWS;
		`,
		tpmTableName,
		tableName,
	)

	_, err = tx.Exec(
		context,
		query,
	)
	if err != nil {
		return fmt.Errorf("could not copy into table %v: could not create temporary table %v: %v", tableName, tpmTableName, err)
	}

	_, err = tx.CopyFrom(
		context,
		pgx.Identifier{tpmTableName},
		columnNames,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("could not copy into table %v: could not copy into temporary table %v: %v", tableName, tpmTableName, err)
	}

	for i, column := range updateColumns {
		updateColumns[i] = fmt.Sprintf("%v = EXCLUDED.%v", column, column)
	}
	updateString := strings.Join(updateColumns, ", ")

	query = fmt.Sprintf(
		`
		INSERT INTO %s
			SELECT *
			FROM %s
		ON CONFLICT (id)
		DO UPDATE
			SET %s;
		`,
		tableName,
		tpmTableName,
		updateString,
	)

	_, err = tx.Exec(
		context,
		query,
	)
	if err != nil {
		return fmt.Errorf("could not copy into table %v: could not update with values from temporary table %v: %v", tableName, tpmTableName, err)
	}

	return nil
}

package util

import (
	"context"
	"fmt"
	"luna-backend/crypto"
	"luna-backend/errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

func CopyAndUpdate(tx pgx.Tx, context context.Context, tableName string, columnNames []string, updateColumns []string, rows [][]any) *errors.ErrorTrace {
	randomNumber, tr := crypto.GenerateRandomNumber()
	if tr != nil {
		return tr.
			Append(errors.LvlDebug, "Could not copy into table %v", tableName).
			Append(errors.LvlPlain, "Database error")
	}

	if !isSafe(tableName) {
		// server error, because we never let the user decide the table name in the first place
		return errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Table name %v failed vaildation", tableName).
			Append(errors.LvlDebug, "Could not copy into table %v", tableName).
			Append(errors.LvlPlain, "Database error")
	}
	tmpTableName := fmt.Sprintf("temp_%v_%v", tableName, randomNumber)

	query := fmt.Sprintf(
		`
		CREATE TEMP TABLE %s (
			LIKE %s INCLUDING ALL	
		)	ON COMMIT DELETE ROWS;
		`,
		tmpTableName,
		tableName,
	)

	_, err := tx.Exec(
		context,
		query,
	)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not create temporary table %v", tmpTableName).
			Append(errors.LvlDebug, "Could not copy into table %v", tableName).
			Append(errors.LvlPlain, "Database error")
	}

	_, err = tx.CopyFrom(
		context,
		pgx.Identifier{tmpTableName},
		columnNames,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not copy into temporary table %v", tmpTableName).
			Append(errors.LvlDebug, "Could not copy into table %v", tableName).
			Append(errors.LvlPlain, "Database error")
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
		tmpTableName,
		updateString,
	)

	_, err = tx.Exec(
		context,
		query,
	)
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not update with values from temporary table %v", tmpTableName).
			Append(errors.LvlDebug, "Could not copy into table %v", tableName).
			Append(errors.LvlPlain, "Database error")
	}

	return nil
}

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

func CopyAndUpdate(
	tx pgx.Tx,
	context context.Context,
	tableName string,
	idColumnName string,
	columnNames []string,
	updateColumns []string,
	rows [][]any,
	deleteUnknown bool,
	deleteCondition string,
	deleteArgument any,
	hasDisplayOrder bool,
	displayOrderGroupByColumn string,
	displayOrderDefaultSortColumn string,
) *errors.ErrorTrace {
	randomNumber, tr := crypto.GenerateRandomNumber()
	if tr != nil {
		return tr.
			Append(errors.LvlDebug, "Could not copy into table %v", tableName).
			Append(errors.LvlPlain, "Database error")
	}

	// TODO: maybe add the same check for idColumnName and columnNames
	if !isSafe(tableName) {
		// server error, because we never let the user decide the table name in the first place
		return errors.New().Status(http.StatusInternalServerError).
			Append(errors.LvlDebug, "Table name %v failed vaildation", tableName).
			Append(errors.LvlDebug, "Could not copy into table %v", tableName).
			Append(errors.LvlPlain, "Database error")
	}
	tmpTableName := fmt.Sprintf("temp_%v_%v", tableName, randomNumber)

	// Create a temporary table for the new elements
	var query string
	if hasDisplayOrder {
		query = fmt.Sprintf(
			`
			CREATE TEMP TABLE %s (
				LIKE %s
				INCLUDING ALL	
				EXCLUDING CONSTRAINTS
			)	ON COMMIT DELETE ROWS;
			`,
			tmpTableName,
			tableName,
		)
	} else {
		query = fmt.Sprintf(
			`
			CREATE TEMP TABLE %s (
				LIKE %s
				INCLUDING ALL	
			)	ON COMMIT DELETE ROWS;
			`,
			tmpTableName,
			tableName,
		)
	}

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

	// Insert all the new elements into the created table
	if hasDisplayOrder {
		for i, row := range rows {
			rows[i] = append(row, "0")
		}

		_, err = tx.CopyFrom(
			context,
			pgx.Identifier{tmpTableName},
			append(columnNames, "display_order"),
			pgx.CopyFromRows(rows),
		)
	} else {

	}
	if err != nil {
		return errors.New().Status(http.StatusInternalServerError).
			AddErr(errors.LvlDebug, err).
			Append(errors.LvlDebug, "Could not copy into temporary table %v", tmpTableName).
			Append(errors.LvlDebug, "Could not copy into table %v", tableName).
			Append(errors.LvlPlain, "Database error")
	}

	// Delete stale entries
	if deleteUnknown {
		query = fmt.Sprintf(
			`
			DELETE FROM %[1]s
			WHERE %[3]s IN (
				SELECT original.%[3]s
				FROM %[1]s original
				LEFT JOIN %[2]s temporary ON original.%[3]s = temporary.%[3]s
				WHERE temporary.%[3]s IS NULL
				AND %[4]s
			);
			`,
			tableName,
			tmpTableName,
			idColumnName,
			deleteCondition,
		)

		_, err = tx.Exec(
			context,
			query,
			deleteArgument,
		)
		if err != nil {
			return errors.New().Status(http.StatusInternalServerError).
				AddErr(errors.LvlDebug, err).
				Append(errors.LvlDebug, "Could not delete stale entries from table %v", tableName).
				Append(errors.LvlDebug, "Could not copy into table %v", tableName).
				Append(errors.LvlPlain, "Database error")
		}
	}

	// Merge the temporary table with the real one
	for i, column := range updateColumns {
		updateColumns[i] = fmt.Sprintf("%v = EXCLUDED.%v", column, column)
	}
	updateString := strings.Join(updateColumns, ", ")

	if hasDisplayOrder {
		query = fmt.Sprintf(
			`
			WITH existingEntries AS (
				SELECT DISTINCT %[3]s, %[5]s	
				FROM %[1]s
			), newEntries AS (
				SELECT temporary.*, ROW_NUMBER() OVER (
					PARTITION BY temporary.%[5]s
					ORDER BY temporary.%[6]s
				) - 1 as newIndex
				FROM %[2]s temporary
				LEFT JOIN existingEntries ON temporary.%[3]s = existingEntries.%[3]s
				WHERE existingEntries.%[3]s IS NULL
			), maxima AS (
				SELECT %[5]s, MAX(display_order) as maxOrder
				FROM %[1]s
				GROUP BY %[5]s
			), orderedEntries AS (
				SELECT allEntries.*, COALESCE(maxima.maxOrder, 0) + COALESCE(newEntries.newIndex, 0) as new_display_order
				FROM %[2]s allEntries
				LEFT JOIN newEntries ON newEntries.%[3]s = allEntries.%[3]s
				LEFT JOIN maxima ON newEntries.%[5]s = maxima.%[5]s
			)
			INSERT INTO %[1]s (%[7]s, display_order)
				SELECT %[7]s, new_display_order
				FROM orderedEntries
			ON CONFLICT (%[3]s)
			DO UPDATE
				SET %[4]s;
			`,
			tableName,
			tmpTableName,
			idColumnName,
			updateString,
			displayOrderGroupByColumn,
			displayOrderDefaultSortColumn,
			strings.Join(columnNames, ", "),
		)
	} else {
		query = fmt.Sprintf(
			`
			INSERT INTO %s
				SELECT *
				FROM %s
			ON CONFLICT (%s)
			DO UPDATE
				SET %s;
			`,
			tableName,
			tmpTableName,
			idColumnName,
			updateString,
		)
	}

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

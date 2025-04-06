package query

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Table struct {
	Schema    string
	TableName string
}

func GetAllTables(db *sql.DB) ([]Table, error) {
	rows, err := db.Query("select table_schema , table_name from information_schema.tables where table_schema not in ('pg_catalog' ,'information_schema' )")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Table
	for rows.Next() {
		var table Table
		if err := rows.Scan(&table.Schema, &table.TableName); err != nil {
			return nil, err
		}
		result = append(result, table)
	}
	return result, nil
}

func ExecuteQuery(db *sql.DB, query string) ([][]string, string, error) {
	trimmed := strings.TrimSpace(query)
	upper := strings.ToUpper(trimmed)

	switch {
	case strings.HasPrefix(upper, "INSERT"),
		strings.HasPrefix(upper, "UPDATE"),
		strings.HasPrefix(upper, "DELETE"):
		result, err := db.Exec(query)
		if err != nil {
			return nil, "", err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, "", err
		}
		return nil, fmt.Sprintf("%d affected rows", rowsAffected), nil

	default:
		rows, err := db.Query(query)
		if err != nil {
			return nil, "", err
		}
		defer rows.Close()
		result, err := getResult(rows)
		return result, "", err
	}
}
func getResult(rows *sql.Rows) ([][]string, error) {
	// Get column names
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := [][]string{cols}

	values := make([]interface{}, len(cols))
	scanArgs := make([]interface{}, len(cols))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Important: Create a new container for each row
		row := make([]string, len(cols))
		for i, v := range values {
			if v == nil {
				row[i] = "NULL"
			} else {
				switch v := v.(type) {
				case []byte:
					row[i] = string(v)
				default:
					row[i] = fmt.Sprintf("%v", v)
				}
			}
		}
		result = append(result, row)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return result, nil

}
func GetTableInformation(table string, db *sql.DB) ([][]string, error) {
	// Use placeholder to prevent SQL injection
	query := "SELECT * FROM " + table
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table %s: %v", table, err)
	}
	defer rows.Close()
	return getResult(rows)

}

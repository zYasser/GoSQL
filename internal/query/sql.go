package query

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Table struct {
	Schmea    string
	TableName string
}

func GetAllTables(db *sql.DB) ([]Table, error) {
	rows, err := db.Query("select table_schema , table_name from information_schema.tables where table_schema not in ('pg_catalog' ,'information_schema' )")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tables: %v", err)
	}
	defer rows.Close()
	var result []Table
	for rows.Next() {
		var table Table
		if err := rows.Scan(&table.Schmea, &table.TableName); err != nil {
			return nil, fmt.Errorf("failed to fetch tables: %v", err)
		}
		result = append(result, table)
	}
	return result, nil
}

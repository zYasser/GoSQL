package services

import (
	"GoSQL/internal/config"
	"GoSQL/internal/query"
	"context"
	"fmt"
	"strings"
)

type UpdateQueryParams struct {
	Key    string
	Values map[string]string
}

func GetTables(ctx context.Context) (map[string][]string, error) {
	db := ctx.Value("db").(*config.DbConfig)
	if db == nil {
		return nil, fmt.Errorf("make sure to chose database profile")
	}
	result, err := query.GetAllTables(db.Connection)
	if err != nil {
		return nil, err
	}
	rows := make(map[string][]string)
	for _, item := range result {
		rows[item.Schema] = append(rows[item.Schema], item.TableName)
	}
	return rows, nil

}

func GenerateSQL(updates []UpdateQueryParams, pk string, table string, ctx context.Context) error {
	db := ctx.Value("db").(*config.DbConfig)
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	for _, update := range updates {
		setParts := []string{}
		for col, val := range update.Values {
			setParts = append(setParts, fmt.Sprintf(`"%s" = '%s'`, col, val))
		}
		queryString := fmt.Sprintf(`UPDATE "%s" SET %s WHERE "%s" = '%s';`,
			table,
			strings.Join(setParts, ", "),
			pk,
			update.Key,
		)
		_, _, err := query.ExecuteQuery(tx, queryString)
		if err != nil {
			tx.Rollback()
			return err
		}

	}
	err = tx.Commit()
	if err != nil {
		return err

	}
	tx.Commit()
	return nil

}

func FetchTableData(ctx context.Context, tableName string) ([][]string, string, error) {
	db := ctx.Value("db").(*config.DbConfig)
	result, pk, err := query.GetTableInformation(tableName, db.Connection)
	if err != nil {
		return nil, "", err
	}
	return result, pk, nil

}

func ExecuteQuery(ctx context.Context, queryString string) ([][]string, string, error) {

	db := ctx.Value("db").(*config.DbConfig)
	tx, err := db.Connection.Begin()
	if err != nil {
		return nil, "", err

	}

	result, message, err := query.ExecuteQuery(tx, queryString)
	if err != nil {
		tx.Rollback()

		return nil, "", err
	}
	err = tx.Commit()
	if err != nil {
		return nil, "", err

	}

	return result, message, nil

}

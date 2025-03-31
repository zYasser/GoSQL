package services

import (
	"GoSQL/internal/config"
	"GoSQL/internal/query"
	"context"
)

func GetTables(ctx context.Context) (map[string][]string, error) {
	db := ctx.Value("db").(*config.DbConfig)
		result, err := query.GetAllTables(db.Connection)
	if err != nil {
		return nil, err
	}
	rows := make(map[string][]string)
	for _, item := range result {
		rows[item.Schmea] = append(rows[item.Schmea], item.TableName)
	}
	return rows, nil

}

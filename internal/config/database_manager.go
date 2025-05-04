package config

import (
	"GoSQL/helpers"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseConnectionInput struct {
	ID           string
	ProfileName  string `json:"ProfileName"`
	Host         string `json:"Host"`
	Port         string `json:"Port"`
	Username     string `json:"Username"`
	Password     string `json:"Password"`
	DatabaseName string `json:"DatabaseName"`
}

type DbConfig struct {
	Connection *sql.DB
}

func ConnectToDb(d DatabaseConnectionInput, ctx context.Context, encrypted bool) error {
	if encrypted {
		originalPassword, err := helpers.DecryptAES(d.Password)
		if err != nil {
			return err
		}
		d.Password = originalPassword
	}
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		d.Username, d.Password, d.Host, d.Port, d.DatabaseName)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}
	dbConfig := ctx.Value("db").(*DbConfig)
	dbConfig.Connection = db

	return nil

}

package config

import (
	"GoSQL/internal/constants"
	"database/sql"
	"fmt"
)

type DbConfig struct {
	Database   constants.Database
	Connection *sql.DB
}

type DatabaseConfig struct {
	ProfileName  string
	DatabaseType constants.Database
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

func NewDb(userInput string) (constants.Database, error) {
	db, err := constants.MapUserInputToDatabase(userInput)
	if err != nil {
		fmt.Printf("Fail To Map input to database ")
	}
	return db, nil
}

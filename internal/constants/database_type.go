package constants

import (
	"errors"

	_ "github.com/lib/pq"
)

type Database string

type Db struct {
	Database Database
}

const (
	Postgres Database = "postgres"
	MySQL    Database = "mysql"
	SQLite   Database = "sqlite"
)

func GetAllDatabaseType() []string {
	return []string{string(Postgres), string(MySQL), string(SQLite)}
}


func MapUserInputToDatabase(userInput string) (Database, error) {
	switch userInput {
	case "1":
		return Postgres, nil
	case "2":
		return MySQL, nil
	case "3":
		return SQLite, nil
	default:
		return "", errors.New("Invalid Database Choose from below:")
	}
}

package config

import "errors"

type Database string

type Db struct {
	Database Database
}

const (
	Postgres Database = "postgres"
	MySQL    Database = "mysql"
	SQLite   Database = "sqlite"
)

// NewDb initializes a new Db struct with the selected database
func NewDb(userInput string) (*Db, error) {

	db, err := mapUserInputToDatabase(userInput)
	if err != nil {
		return nil, err
	}
	return &Db{Database: db}, nil
}

// mapUserInputToDatabase converts user input to a database type
func mapUserInputToDatabase(userInput string) (Database, error) {
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

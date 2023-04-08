package database

import (
	"github.com/surrealdb/surrealdb.go"
)

type Database struct {
	Client *surrealdb.DB
}

func NewDatabase(databaseURL string) (*Database, error) {
	client, err := surrealdb.New(databaseURL)
	if err != nil {
		return nil, err
	}

	return &Database{
		Client: client,
	}, nil
}

func (db *Database) Init(username string, password string) error {
	details := map[string]any{
		"user": username,
		"pass": password,
	}

	if _, err := db.Client.Signin(details); err != nil {
		return err
	}

	if _, err := db.Client.Use("line-chatgpt", "line-chatgpt"); err != nil {
		return err
	}

	return nil
}

type Response struct {
	Result []any  `json:"result"`
	Status string `json:"status"`
	Time   string `json:"time"`
}

type Person struct {
	Age  int    `json:"age"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Sex  bool   `json:"sex"`
}

type QueryResult struct {
	Result []map[string]interface{}
}

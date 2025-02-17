package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDb() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "store")

	return database, err
}

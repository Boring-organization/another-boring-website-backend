package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func InitDb() *sql.DB {
	database, err := sql.Open("sqlite3", "store.db")

	if err != nil {
		panic(err)
	}

	return database
}

func (db *Database) CloseDb() {
	err := db.DB.Close()

	if err != nil {
		panic(err)
	}
}

func (db *Database) ExecuteOperation(query string, args ...any) sql.Result {
	result, err := db.DB.Exec(query, args...)

	if err != nil {
		panic(err)
	}

	return result
}

func (db *Database) Query(query string, args ...any) *sql.Rows {
	rows, err := db.DB.Query(query, args...)

	if err != nil {
		panic(err)
	}

	return rows
}

func (db *Database) QueryRow(query string, args ...any) *sql.Row {
	row := db.DB.QueryRow(query, args...)

	return row
}

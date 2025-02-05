package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func InitDb() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "store.db")

	return database, err
}

func (db *Database) CloseDb() {
	err := db.DB.Close()

	if err != nil {
		panic(err)
	}
}

func (db *Database) ExecuteOperation(query string, args ...any) (sql.Result, error) {
	result, err := db.DB.Exec(query, args...)

	return result, err
}

func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	rows, err := db.DB.Query(query, args...)

	return rows, err
}

func (db *Database) QueryRow(query string, args ...any) *sql.Row {
	row := db.DB.QueryRow(query, args...)

	return row
}

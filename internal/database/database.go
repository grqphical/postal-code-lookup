package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// initalizes the database connection
func InitDB() (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", "database.sqlite3")
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS Municipalities (id INTEGER PRIMARY KEY, fsa TEXT, municipality TEXT);")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

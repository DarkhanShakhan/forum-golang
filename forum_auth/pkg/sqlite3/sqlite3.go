package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./session.db")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	sessions := `
	CREATE TABLE IF NOT EXISTS sessions (
		user_id INTEGER NOT NULL UNIQUE,
		token TEXT NOT NULL UNIQUE,
		expiry_date TEXT NOT NULL
	);`
	_, err = db.Exec(sessions)
	if err != nil {
		return nil, err
	}
	return db, nil
}

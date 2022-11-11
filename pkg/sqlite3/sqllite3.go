package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./forum.db")

	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys = ON;")
	users := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT,
		registration_date TEXT
		);
	`
	_, err = db.Exec(users)
	posts := `
	CREATE TABLE posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		date TEXT,
		title TEXT,
		content TEXT
		);
	`
	_, err = db.Exec(posts)

	return db, nil
}

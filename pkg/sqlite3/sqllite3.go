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

	users := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT,
		registration_date TEXT,
		);
	`
	_, err = db.Exec(users)

	query := `INSERT INTO users(name, email, password, registration_date) VALUES("user1", "user1@mail.com", "password", "03-10-2022")`
	_, err = db.Exec(query)
	return db, nil
}

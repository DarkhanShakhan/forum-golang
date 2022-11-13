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
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT,
		registration_date TEXT
		);
	`
	_, err = db.Exec(users)
	// db.Exec(`INSERT INTO users(name,email, password) VALUES ("user1", "user1", "user1");`)
	// db.Exec(`INSERT INTO users(name, email, password) VALUES ("user2", "user2", "user2");`)
	// db.Exec(`INSERT INTO users(name, email, password) VALUES ("user3", "user3", "user3");`)
	posts := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		date TEXT,
		title TEXT,
		content TEXT
		);
	`
	_, err = db.Exec(posts)
	postReactions := `
	CREATE TABLE IF NOT EXISTS post_reactions (
		post_id INTEGER,
		user_id INTEGER,
		date TEXT,
		like INTEGER,
		PRIMARY KEY(post_id, user_id)
		);
	`
	_, err = db.Exec(postReactions)
	return db, nil
}

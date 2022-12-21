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
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys = ON;")
	users := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL UNIQUE,
		registration_date TEXT
		);
	`
	_, err = db.Exec(users)
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return nil, err
	}
	postReactions := `
	CREATE TABLE IF NOT EXISTS post_reactions (
		post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		date TEXT,
		like INTEGER,
		UNIQUE(post_id, user_id)
		);
	`
	_, err = db.Exec(postReactions)
	if err != nil {
		return nil, err
	}
	categories := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT UNIQUE
	);`
	_, err = db.Exec(categories)
	if err != nil {
		return nil, err
	}
	// FIXME: add categories
	postCategories := `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
		category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
		UNIQUE(post_id, category_id)
	);`
	_, err = db.Exec(postCategories)
	if err != nil {
		return nil, err
	}
	comments := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		date STRING,
		content STRING
	);`
	_, err = db.Exec(comments)
	if err != nil {
		return nil, err
	}
	commentReactions := `
	CREATE TABLE IF NOT EXISTS comment_reactions (
		comment_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		date TEXT,
		like INTEGER,
		UNIQUE(comment_id, user_id)
		);
	`
	_, err = db.Exec(commentReactions)
	if err != nil {
		return nil, err
	}
	return db, nil
}

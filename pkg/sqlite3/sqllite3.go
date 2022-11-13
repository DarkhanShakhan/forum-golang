package sqlite3

import (
	"database/sql"
	"fmt"

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
		post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		date TEXT,
		like INTEGER,
		UNIQUE(post_id, user_id)
		);
	`
	_, err = db.Exec(postReactions)
	fmt.Println(err)
	categories := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT UNIQUE
	);`
	db.Exec(categories)
	cat1 := `
	Insert into categories(title) values("sql");`
	cat2 := `insert into categories(title) values("python");`
	cat3 := `insert into categories(title) values("golang");`
	db.Exec(cat1)
	db.Exec(cat2)
	db.Exec(cat3)
	postCategories := `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
		category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
		UNIQUE(post_id, category_id)
	);`
	db.Exec(postCategories)
	return db, nil
}

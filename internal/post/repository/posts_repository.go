package repository

import (
	"database/sql"
	"forum/internal/entity"
)

const (
	SELECT_QUERY = "SELECT * FROM"
	DELETE_QUERY = "DELETE FROM"
	POSTS        = " posts"
	BY           = " WHERE"
	BY_ID        = BY + " id = ?"
	BY_USER_ID   = BY + " user_id = ?"
)

type PostsRepository struct {
	db *sql.DB
}

func NewPostsRepository(db *sql.DB) *PostsRepository {
	return &PostsRepository{db}
}

func (pr *PostsRepository) FetchById(id int) (entity.Post, error) {
	post := entity.Post{}
	tx, err := pr.db.Begin()
	if err != nil {
		return post, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(SELECT_QUERY + POSTS + BY_ID)
	if err != nil {
		return post, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return post, err
	}
	if rows.Next() {
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
	}
	if err = tx.Commit(); err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (pr *PostsRepository) FetchByUserId(id int) ([]entity.Post, error) {
	posts := []entity.Post{}
	tx, err := pr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(SELECT_QUERY + POSTS + BY_USER_ID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := entity.Post{}
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
		posts = append(posts, post)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostsRepository) FetchByCategoryId(id int) ([]entity.Post, error) {
	posts := []entity.Post{}
	tx, err := pr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("SELECT post_id FROM post_categories WHERE category_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := entity.Post{}
		rows.Scan(&post.Id)
		posts = append(posts, post)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostsRepository) Store(post entity.Post) (int64, error) {
	tx, err := pr.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO posts(user_id, date, title, content) VALUES(?,?,?,?);`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(post.User.Id, post.Date, post.Title, post.Content)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

//for future use
// func (pr *PostsRepository) Update(entity.Post) error {
// 	return nil
// }

// func (pr *PostsRepository) Delete(id int) error {
// 	result, err := pr.db.Exec(DELETE_QUERY+POSTS+BY_ID, id)
// 	if err != nil {
// 		return err
// 	}
// 	nbr, err := result.RowsAffected()
// 	if nbr > 1 {
// 		return errors.New("more than one row has been affected")
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

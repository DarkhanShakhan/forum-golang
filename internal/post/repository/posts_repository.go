package repository

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

const (
	SELECT_QUERY = "SELECT * FROM"
	DELETE_QUERY = "DELETE FROM"
	POSTS        = " posts"
	BY_ID        = " WHERE id = ?"
	BY_USER_ID   = " WHERE user_id = ?"
)

type PostsRepository struct {
	db *sql.DB
}

func NewPostsRepository(db *sql.DB) *PostsRepository {
	return &PostsRepository{db}
}

func (pr *PostsRepository) FetchById(id int) (entity.Post, error) {
	post := entity.Post{}
	rows, err := pr.db.Query(SELECT_QUERY+POSTS+BY_ID, id)
	if err != nil {
		return post, err
	}
	if rows.Next() {
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
	}
	return post, nil
}

func (pr *PostsRepository) FetchByUserId(id int) ([]entity.Post, error) {
	posts := []entity.Post{}
	rows, err := pr.db.Query(SELECT_QUERY+POSTS+BY_USER_ID, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := entity.Post{}
		rows.Scan(&post.Id, &post.User.Id, &post.Date, &post.Title, &post.Content)
		posts = append(posts, post)
	}
	return posts, nil
}

func (pr *PostsRepository) Store(post entity.Post) (int64, error) {
	res, err := pr.db.Exec(`INSERT INTO posts(user_id, date, title, content) VALUES(?,?,?,?);`, post.User.Id, post.Date, post.Title, post.Content)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

//for future use
func (pr *PostsRepository) Update(entity.Post) error {
	return nil
}

func (pr *PostsRepository) Delete(id int) error {
	result, err := pr.db.Exec(DELETE_QUERY+POSTS+BY_ID, id)
	if err != nil {
		return err
	}
	nbr, err := result.RowsAffected()
	if nbr > 1 {
		return errors.New("more than one row has been affected")
	}
	if err != nil {
		return err
	}
	return nil
}

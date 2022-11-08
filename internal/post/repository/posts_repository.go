package repository

import (
	"database/sql"
	"forum/internal/entity"
)

const (
	SELECT_QUERY = "SELECT * FROM"
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

func (pr *PostsRepository) FetchById(int) (entity.Post, error) {
	return entity.Post{}, nil
}

func (pr *PostsRepository) FetchByUserId(int) ([]entity.Post, error) {
	return nil, nil
}

func (pr *PostsRepository) Store(entity.Post) (int, error) {
	return 0, nil
}

func (pr *PostsRepository) Update(entity.Post) error {
	return nil
}

func (pr *PostsRepository) Delete(int) error {
	return nil
}

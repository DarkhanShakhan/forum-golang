package repository

import (
	"database/sql"
	"forum/internal/entity"
)

type CommentsRepository struct {
	db *sql.DB
}

func NewCommentsRepository(db *sql.DB) *CommentsRepository {
	return &CommentsRepository{db}
}

func (cr *CommentsRepository) FetchById(id int) (entity.Comment, error) {
	return entity.Comment{}, nil
}

func (cr *CommentsRepository) FetchByPostId(id int) ([]entity.Comment, error) {
	return nil, nil
}

func (cr *CommentsRepository) FetchByUserId(id int) ([]entity.Comment, error) {
	return nil, nil
}

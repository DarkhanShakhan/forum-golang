package repository

import (
	"database/sql"
	"forum/internal/entity"
)

type CommentReactionsRepository struct {
	db *sql.DB
}

func NewCommentReactionsRepository(db *sql.DB) *CommentReactionsRepository {
	return &CommentReactionsRepository{db}
}

func (crr *CommentReactionsRepository) FetchByCommentId(id int, like bool) ([]entity.Reaction, error) {
	return nil, nil
}

func (crr *CommentReactionsRepository) FetchByUserId(id int, like bool) ([]entity.Reaction, error) {
	return nil, nil
}

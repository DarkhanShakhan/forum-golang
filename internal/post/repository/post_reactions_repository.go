package repository

import (
	"database/sql"
	"forum/internal/entity"
)

const (
	POST_REACTIONS   = " post_reactions"
	BY_POST_AND_LIKE = " WHERE post_id = ? AND like = ?"
)

type PostReactionsRepository struct {
	db *sql.DB
}

func NewPostReactionsRepository(db *sql.DB) *PostReactionsRepository {
	return &PostReactionsRepository{db}
}

func (rr *PostReactionsRepository) FetchByPostId(int, bool) ([]entity.Reaction, error) {
	return nil, nil
}

func (rr *PostReactionsRepository) FetchByUserId(int) ([]entity.PostReaction, error) {
	return nil, nil
}

package app

import "forum_app/internal/entity"

type UserUsecase interface {
	FetchById(int) (entity.User, error)
	FetchAll() ([]entity.User, error)
	FetchByEmail(string) (entity.User, error)
}

type PostUsecase interface {
	FetchById(int) (entity.Post, error)
	FetchAll() ([]entity.Post, error)
	FetchCategoryPosts(entity.Category) (entity.Category, error)
	Store(entity.Post) (int64, error)
	StorePostReaction(entity.PostReaction) error
	UpdatePostReaction(entity.PostReaction) error
	DeletePostReaction(entity.PostReaction) error
}

type CommentUsecase interface {
	FetchById(int) (entity.Comment, error)
	Store(entity.Comment) (int64, error)
	StoreCommentReaction(entity.CommentReaction) error
	UpdateCommentReaction(entity.CommentReaction) error
	DeleteCommentReaction(entity.CommentReaction) error
}

package app

import (
	"context"
	"forum_app/internal/entity"
)

type UserUsecase interface {
	FetchById(context.Context, int) (entity.User, error)
	FetchAll(context.Context) ([]entity.User, error)
	FetchByEmail(context.Context, string) (entity.User, error)
}

type PostUsecase interface {
	FetchById(context.Context, int) (entity.Post, error)
	FetchAll(context.Context) ([]entity.Post, error)
	FetchCategoryPosts(context.Context, int) (entity.Category, error)
	Store(context.Context, entity.Post) (int64, error)
	StorePostReaction(context.Context, entity.PostReaction) error
	UpdatePostReaction(context.Context, entity.PostReaction, chan error)
	DeletePostReaction(context.Context, entity.PostReaction) error
}

type CommentUsecase interface {
	FetchById(context.Context, int) (entity.Comment, error)
	Store(context.Context, entity.Comment) (int64, error)
	StoreCommentReaction(context.Context, entity.CommentReaction) error
	UpdateCommentReaction(context.Context, entity.CommentReaction) error
	DeleteCommentReaction(context.Context, entity.CommentReaction) error
}

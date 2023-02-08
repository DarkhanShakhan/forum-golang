package usecase

import (
	"context"
	"forum_app/internal/entity"
)

type CommentsRepository interface {
	FetchById(context.Context, int) (entity.Comment, error)
	Store(context.Context, entity.Comment) (int64, error)
}

type CommentReactionsRepository interface {
	FetchByCommentId(context.Context, int, bool) ([]entity.Reaction, error)
	StoreReaction(context.Context, entity.CommentReaction) error
	UpdateReaction(context.Context, entity.CommentReaction) error
	DeleteReaction(context.Context, entity.CommentReaction) error
}

type PostsRepository interface {
	FetchById(context.Context, int) (entity.Post, error)
}

type UsersRepository interface {
	FetchById(context.Context, int) (entity.User, error)
}

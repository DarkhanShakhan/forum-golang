package usecase

import (
	"context"
	"forum_app/internal/entity"
)

type UsersRepository interface {
	FetchById(context.Context, int) (entity.User, error)
	FetchAll(context.Context) ([]entity.User, error)
	FetchByEmail(context.Context, string) (entity.User, error)
	// Update(entity.User) error
	// DeleteById(int) error
}

type PostsRepository interface {
	FetchByUserId(context.Context, int) ([]entity.Post, error)
}

type PostReactionsRepository interface {
	FetchByUserId(context.Context, int, bool) ([]entity.PostReaction, error)
}

type CommentReactionsRepository interface {
	FetchByUserId(context.Context, int, bool) ([]entity.CommentReaction, error)
}

type CommentRepository interface {
	FetchByUserId(context.Context, int) ([]entity.Comment, error)
}

package usecase

import (
	"context"
	"forum_app/internal/entity"
)

type PostsRepository interface {
	FetchById(context.Context, int) (entity.Post, error)
	FetchByCategoryId(context.Context, int) ([]entity.Post, error)
	FetchAll(context.Context) ([]entity.Post, error)
	Store(context.Context, entity.Post) (int64, error)
}

type PostReactionsRepository interface {
	FetchByPostId(context.Context, int, bool) ([]entity.Reaction, error)
	StoreReaction(context.Context, entity.PostReaction) error
	UpdateReaction(context.Context, entity.PostReaction) error
	DeleteReaction(context.Context, entity.PostReaction) error
}

type UsersRepository interface {
	FetchById(context.Context, int) (entity.User, error)
}

type CommentsRepository interface {
	FetchByPostId(context.Context, int) ([]entity.Comment, error)
}

type CommentReactionsRepository interface {
	FetchByCommentId(context.Context, int, bool) ([]entity.Reaction, error)
}

type CategoriesRepository interface {
	FetchById(context.Context, int) (entity.Category, error)
	FetchByPostId(context.Context, int) ([]entity.Category, error)
	FetchAllCategories(context.Context) ([]entity.Category, error)
}

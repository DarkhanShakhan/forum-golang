package app

import (
	"context"
	"forum_app/internal/entity"
)

type UserUsecase interface {
	FetchById(context.Context, int, chan entity.UserResult)
	FetchAll(context.Context, chan entity.UsersResult)
	FetchByEmail(context.Context, string, chan entity.UserResult)
	Store(context.Context, entity.User, chan entity.Result)
}

type PostUsecase interface {
	FetchById(context.Context, int, chan entity.PostResult)
	FetchAll(context.Context, chan entity.PostsResult)
	FetchCategoryPosts(context.Context, int, chan entity.CatResult)
	Store(context.Context, entity.Post, chan entity.Result)
	StorePostReaction(context.Context, entity.PostReaction, chan error)
	UpdatePostReaction(context.Context, entity.PostReaction, chan error)
	DeletePostReaction(context.Context, entity.PostReaction, chan error)
}

type CommentUsecase interface {
	// FetchById(context.Context, int, chan entity.CommentResult)
	Store(context.Context, entity.Comment, chan entity.Result)
	StoreCommentReaction(context.Context, entity.CommentReaction, chan error)
	UpdateCommentReaction(context.Context, entity.CommentReaction, chan error)
	DeleteCommentReaction(context.Context, entity.CommentReaction, chan error)
}

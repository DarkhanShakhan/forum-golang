package usecase

import "forum/internal/entity"

type CommentsRepository interface {
	FetchById(int) (entity.Comment, error)
}

type CommentReactionsRepository interface {
	FetchByCommentId(int, bool) ([]entity.Reaction, error)
}

type PostsRepository interface {
	FetchByCommentId(int) (entity.Post, error)
}

type UsersRepository interface {
	FetchByCommentId(int) (entity.User, error)
}

package usecase

import "forum/internal/entity"

type CommentsRepository interface {
	FetchById(int) (entity.Comment, error)
	Store(entity.Comment) (int, error)
	Update(entity.Comment) error
	DeleteById(int) error
}

type CommentReactionsRepository interface {
	FetchByCommentId(int, bool) ([]entity.Reaction, error)
	StoreReaction(entity.CommentReaction) error
	UpdateReaction(entity.CommentReaction) error
	DeleteReaction(entity.CommentReaction) error
}

type PostsRepository interface {
	FetchByCommentId(int) (entity.Post, error)
}

type UsersRepository interface {
	FetchByCommentId(int) (entity.User, error)
}

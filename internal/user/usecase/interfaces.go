package usecase

import "forum/internal/entity"

type UsersRepository interface {
	FetchById(int) (entity.User, error)
	FetchAll() ([]entity.User, error)
	Update(entity.User) error
	Delete(int) error
}

type PostsRepository interface {
	FetchByUserId(int) ([]entity.Post, error)
}

type PostReactionsRepository interface {
	FetchByUserId(int, bool) ([]entity.PostReaction, error)
}

type CommentReactionsRepository interface {
	FetchByUserId(int, bool) ([]entity.CommentReaction, error)
}

type CommentRepository interface {
	FetchByUserId(int) ([]entity.Comment, error)
}

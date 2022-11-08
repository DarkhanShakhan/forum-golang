package usecase

import "forum/internal/entity"

type PostsRepository interface {
	FetchById(int) (entity.Post, error)
	FetchByUserId(int) (entity.Post, error)
	Store(entity.Post) (entity.Post, error)
	Update(entity.Post) error
	Delete(int) error
}

type PostReactionsRepository interface {
	FetchByPostId(int) (entity.PostReaction, error)
}

type UsersRepository interface {
	FetchById(int) (entity.User, error)
}

type CommentsRepository interface {
	FetchByPostId(int) ([]entity.Comment, error)
}

type CategoriesRepository interface {
	FetchByPostId(int) ([]entity.Category, error)
}

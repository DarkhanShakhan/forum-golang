package usecase

import "forum/internal/entity"

type PostsRepository interface {
	FetchById(int) (entity.Post, error)
	Store(entity.Post) (int, error)
	Update(entity.Post) error
	Delete(int) error
}

type PostReactionsRepository interface {
	FetchByPostId(int, bool) ([]entity.Reaction, error)
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

package usecase

import "forum/internal/entity"

type PostsRepository interface {
	FetchById(int) (entity.Post, error)
	FetchByCategoryId(int) ([]entity.Post, error)
	FetchAll() ([]entity.Post, error)
	Store(entity.Post) (int64, error)
	// Update(entity.Post) error
	// Delete(int) error
}

type PostReactionsRepository interface {
	FetchByPostId(int, bool) ([]entity.Reaction, error)
	StoreReaction(entity.PostReaction) error
	UpdateReaction(entity.PostReaction) error
	DeleteReaction(entity.PostReaction) error
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

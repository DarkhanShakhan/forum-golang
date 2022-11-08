package usecase

import (
	"forum/internal/entity"
	"log"
)

type PostsUsecase struct {
	postRepo         PostsRepository
	postReactionRepo PostReactionsRepository
	commentRepo      CommentsRepository
	categoriesRepo   CategoriesRepository
}

func NewPostsUsecase(postRepo PostsRepository, postReactionsRepo PostReactionsRepository, commentRepo CommentsRepository, categoriesRepo CategoriesRepository) *PostsUsecase {
	return &PostsUsecase{postRepo: postRepo, postReactionRepo: postReactionsRepo, commentRepo: commentRepo, categoriesRepo: categoriesRepo}
}

func (u *PostsUsecase) FetchById(id int) (entity.Post, error) {
	post, err := u.postRepo.FetchById(id)
	if err != nil {
		return entity.Post{}, err
	}
	comments := make(chan []entity.Comment)
	categories := make(chan []entity.Category)
	errComments := make(chan error)
	errCategories := make(chan error)
	go u.fetchCategories(id, categories, errCategories)
	go u.fetchComments(id, comments, errComments)
	if err = <-errCategories; err != nil {
		log.Println(err)
	}
	post.Category = <-categories
	if err = <-errComments; err != nil {
		log.Println(err)
	}
	post.Comments = <-comments
	post.CountTotals()
	return post, nil
}

func (u *PostsUsecase) fetchComments(id int, comments chan []entity.Comment, errComments chan error) {
	tempComments, err := u.commentRepo.FetchByPostId(id)
	comments <- tempComments
	errComments <- err
}

func (u *PostsUsecase) fetchCategories(id int, categories chan []entity.Category, errCategories chan error) {
	tempCategories, err := u.categoriesRepo.FetchByPostId(id)
	categories <- tempCategories
	errCategories <- err
}

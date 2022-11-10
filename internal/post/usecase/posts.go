package usecase

import (
	"forum/internal/entity"
	"log"
)

type PostsUsecase struct {
	postsRepo         PostsRepository
	postReactionsRepo PostReactionsRepository
	commentsRepo      CommentsRepository
	categoriesRepo    CategoriesRepository
}

func NewPostsUsecase(postsRepo PostsRepository, postReactionsRepo PostReactionsRepository, commentsRepo CommentsRepository, categoriesRepo CategoriesRepository) *PostsUsecase {
	return &PostsUsecase{postsRepo: postsRepo, postReactionsRepo: postReactionsRepo, commentsRepo: commentsRepo, categoriesRepo: categoriesRepo}
}

func (u *PostsUsecase) FetchById(id int) (entity.Post, error) {
	post, err := u.postsRepo.FetchById(id)
	if err != nil {
		return entity.Post{}, err
	}
	comments := make(chan []entity.Comment)
	categories := make(chan []entity.Category)
	likes := make(chan []entity.Reaction)
	dislikes := make(chan []entity.Reaction)
	errComments := make(chan error)
	errCategories := make(chan error)
	errLikes := make(chan error)
	errDislikes := make(chan error)
	go u.fetchCategories(id, categories, errCategories)
	go u.fetchComments(id, comments, errComments)
	go u.fetchLikes(id, likes, errLikes)
	go u.fetchDislikes(id, dislikes, errDislikes)
	if err = <-errCategories; err != nil {
		log.Println(err)
	}
	post.Category = <-categories
	if err = <-errComments; err != nil {
		log.Println(err)
	}
	post.Comments = <-comments
	if err = <-errLikes; err != nil {
		log.Println(err)
	}
	post.Likes = <-likes
	if err = <-errDislikes; err != nil {
		log.Println(err)
	}
	post.Dislikes = <-dislikes
	post.CountTotals()
	return post, nil
}

func (u *PostsUsecase) fetchComments(id int, comments chan []entity.Comment, errComments chan error) {
	tempComments, err := u.commentsRepo.FetchByPostId(id)
	comments <- tempComments
	errComments <- err
}

func (u *PostsUsecase) fetchCategories(id int, categories chan []entity.Category, errCategories chan error) {
	tempCategories, err := u.categoriesRepo.FetchByPostId(id)
	categories <- tempCategories
	errCategories <- err
}

func (u *PostsUsecase) fetchLikes(id int, likes chan []entity.Reaction, errLikes chan error) {
	tempLikes, err := u.postReactionsRepo.FetchByPostId(id, true)
	likes <- tempLikes
	errLikes <- err
}

func (u *PostsUsecase) fetchDislikes(id int, dislikes chan []entity.Reaction, errDislikes chan error) {
	tempDislikes, err := u.postReactionsRepo.FetchByPostId(id, false)
	dislikes <- tempDislikes
	errDislikes <- err
}

func (u *PostsUsecase) FetchCategoryById(id int) (entity.Category, error) {
	category, err := u.categoriesRepo.FetchById(id)
	if err != nil {
		return entity.Category{}, err
	}
	category.Posts, err = u.postsRepo.FetchByCategory(id)
	if err != nil {
		log.Println(err)
	}
	return category, nil
}

package usecase

import (
	"fmt"
	"forum/internal/entity"
	"log"
)

type PostsUsecase struct {
	postsRepo         PostsRepository
	postReactionsRepo PostReactionsRepository
	commentsRepo      CommentsRepository
	categoriesRepo    CategoriesRepository
	usersRepo         UsersRepository
}

func NewPostsUsecase(postsRepo PostsRepository, postReactionsRepo PostReactionsRepository, commentsRepo CommentsRepository, categoriesRepo CategoriesRepository, usersRepo UsersRepository) *PostsUsecase {
	return &PostsUsecase{postsRepo: postsRepo, postReactionsRepo: postReactionsRepo, commentsRepo: commentsRepo, categoriesRepo: categoriesRepo, usersRepo: usersRepo}
}

func (u *PostsUsecase) FetchById(id int) (entity.Post, error) {
	post, err := u.postsRepo.FetchById(id)
	if err != nil {
		return entity.Post{}, err
	}
	user := make(chan entity.User)
	comments := make(chan []entity.Comment)
	categories := make(chan []entity.Category)
	likes := make(chan []entity.Reaction)
	dislikes := make(chan []entity.Reaction)
	errUser := make(chan error)
	errComments := make(chan error)
	errCategories := make(chan error)
	errLikes := make(chan error)
	errDislikes := make(chan error)
	go u.fetchUser(post.User.Id, user, errUser)
	go u.fetchCategories(id, categories, errCategories)
	go u.fetchComments(id, comments, errComments)
	go u.fetchLikes(id, likes, errLikes)
	go u.fetchDislikes(id, dislikes, errDislikes)
	for i := 0; i < 5; i++ {
		select {
		case post.User = <-user:
			if err = <-errUser; err != nil {
				fmt.Println(err)
			}
		case post.Category = <-categories:
			if err = <-errCategories; err != nil {
				fmt.Println(err)
			}
		case post.Comments = <-comments:
			if err = <-errComments; err != nil {
				fmt.Println(err)
			}
		case post.Likes = <-likes:
			if err = <-errLikes; err != nil {
				fmt.Println(err)
			}
		case post.Dislikes = <-dislikes:
			if err = <-errDislikes; err != nil {
				fmt.Println(err)
			}
		}
	}
	post.CountTotals()
	return post, nil
}
func (u *PostsUsecase) FetchAll() ([]entity.Post, error) {
	return u.postsRepo.FetchAll()
}

func (u *PostsUsecase) fetchUser(id int, user chan entity.User, errUser chan error) {
	tempUser, err := u.usersRepo.FetchById(id)
	user <- tempUser
	errUser <- err
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

func (u *PostsUsecase) FetchCategoryPosts(category entity.Category) (entity.Category, error) {
	var err error
	category.Posts, err = u.postsRepo.FetchByCategoryId(category.Id)
	if err != nil {
		return category, err
	}
	for ix, post := range category.Posts {
		category.Posts[ix], err = u.postsRepo.FetchById(post.Id)
		log.Println(err)
	}
	category.CountTotals()
	return category, nil
}

// func (u *PostsUsecase) FetchAllSorted() ([]entity.Post, error) {
// 	posts, err := u.postsRepo.FetchAllSorted()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return posts, nil
// }

func (u *PostsUsecase) Store(post entity.Post) (int64, error) {
	id, err := u.postsRepo.Store(post)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *PostsUsecase) StorePostReaction(postReaction entity.PostReaction) error {
	return u.postReactionsRepo.StoreReaction(postReaction)
}

func (u *PostsUsecase) UpdatePostReaction(postReaction entity.PostReaction) error {
	return u.postReactionsRepo.UpdateReaction(postReaction)
}

func (u *PostsUsecase) DeletePostReaction(postReaction entity.PostReaction) error {
	return u.postReactionsRepo.DeleteReaction(postReaction)
}

//for future use
// func (u *PostsUsecase) Update(post entity.Post) error {
// 	err := u.postsRepo.Update(post)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (u *PostsUsecase) DeleteById(id int) error {
// 	err := u.postsRepo.Delete(id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

package usecase

import (
	"forum/internal/entity"
	"log"
)

type UsersUsecase struct {
	userRepo             UsersRepository
	postRepo             PostsRepository
	postReactionsRepo    PostReactionsRepository
	commentRepo          CommentRepository
	commentReactionsRepo CommentReactionsRepository
}

func NewUsersUsecase(userRepo UsersRepository, postRepo PostsRepository, postReactionsRepo PostReactionsRepository, commentRepo CommentRepository, commentReactionsRepo CommentReactionsRepository) *UsersUsecase {
	return &UsersUsecase{userRepo: userRepo, postRepo: postRepo, postReactionsRepo: postReactionsRepo, commentRepo: commentRepo, commentReactionsRepo: commentReactionsRepo}
}

func (u *UsersUsecase) FetchById(id int) (entity.User, error) {
	user, err := u.userRepo.FetchById(id)
	if err != nil {
		return entity.User{}, err
	}
	posts := make(chan []entity.Post)
	comments := make(chan []entity.Comment)
	postReactions := make(chan []entity.PostReaction)
	commentReactions := make(chan []entity.CommentReaction)
	errPosts := make(chan error)
	errComments := make(chan error)
	errPostReactions := make(chan error)
	errCommentReactions := make(chan error)
	go u.fetchPosts(id, posts, errPosts)
	go u.fetchComments(id, comments, errComments)
	go u.fetchPostReactions(id, postReactions, errPostReactions)
	go u.fetchCommentReactions(id, commentReactions, errCommentReactions)

	for i := 0; i < 4; i++ {
		select {
		case user.Posts = <-posts:
			if err = <-errPosts; err != nil {
				log.Println(err)
			}
		case user.Comments = <-comments:

			if err = <-errComments; err != nil {
				log.Println(err)
			}
		case user.PostReactions = <-postReactions:
			if err = <-errPostReactions; err != nil {
				log.Println(err)
			}
		case user.CommentReactions = <-commentReactions:
			if err = <-errCommentReactions; err != nil {
				log.Println(err)
			}
		}
	}
	user.CountTotals()
	return user, nil
}

func (u *UsersUsecase) FetchAll() ([]entity.User, error) {
	users, err := u.userRepo.FetchAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UsersUsecase) Update(user entity.User) error {
	return u.userRepo.Update(user)
}

func (u *UsersUsecase) Delete(id int) error {
	return u.userRepo.Delete(id)
}

func (u *UsersUsecase) fetchPosts(id int, posts chan []entity.Post, errPosts chan error) {
	tempPosts, err := u.postRepo.FetchByUserId(id)
	posts <- tempPosts
	errPosts <- err
}

func (u *UsersUsecase) fetchComments(id int, comments chan []entity.Comment, errComments chan error) {
	tempComments, err := u.commentRepo.FetchByUserId(id)
	comments <- tempComments
	errComments <- err
}

func (u *UsersUsecase) fetchPostReactions(id int, postReactions chan []entity.PostReaction, errPostReactions chan error) {
	tempPostReactions, err := u.postReactionsRepo.FetchByUserId(id)
	postReactions <- tempPostReactions
	errPostReactions <- err
}

func (u *UsersUsecase) fetchCommentReactions(id int, commentReactions chan []entity.CommentReaction, errCommentReactions chan error) {
	tempCommentReactions, err := u.commentReactionsRepo.FetchByUserId(id)
	commentReactions <- tempCommentReactions
	errCommentReactions <- err
}

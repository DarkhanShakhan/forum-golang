package usecase

import (
	"forum/internal/forum_app/entity"
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
	postLikes := make(chan []entity.PostReaction)
	postDislikes := make(chan []entity.PostReaction)
	commentLikes := make(chan []entity.CommentReaction)
	commentDislikes := make(chan []entity.CommentReaction)
	errPosts := make(chan error)
	errComments := make(chan error)
	errPostLikes := make(chan error)
	errPostDislikes := make(chan error)
	errCommentLikes := make(chan error)
	errCommentDislikes := make(chan error)
	go u.fetchPosts(id, posts, errPosts)
	go u.fetchComments(id, comments, errComments)
	go u.fetchPostReactions(id, true, postLikes, errPostLikes)
	go u.fetchPostReactions(id, false, postDislikes, errPostDislikes)
	go u.fetchCommentReactions(id, true, commentLikes, errCommentLikes)
	go u.fetchCommentReactions(id, false, commentDislikes, errCommentDislikes)
	for i := 0; i < 6; i++ {
		select {
		case user.Posts = <-posts:
			if err = <-errPosts; err != nil {
				log.Println(err)
			}
		case user.Comments = <-comments:
			if err = <-errComments; err != nil {
				log.Println(err)
			}
		case user.PostLikes = <-postLikes:
			if err = <-errPostLikes; err != nil {
				log.Println(err)
			}
		case user.PostDislikes = <-postDislikes:
			if err = <-errPostDislikes; err != nil {
				log.Println(err)
			}
		case user.CommentLikes = <-commentLikes:
			if err = <-errCommentLikes; err != nil {
				log.Println(err)
			}

		case user.CommentDislikes = <-commentDislikes:
			if err = <-errCommentDislikes; err != nil {
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

// func (u *UsersUsecase) Update(user entity.User) error {
// 	return u.userRepo.Update(user)
// }

// func (u *UsersUsecase) DeleteById(id int) error {
// 	return u.userRepo.DeleteById(id)
// }

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

func (u *UsersUsecase) fetchPostReactions(id int, like bool, postReactions chan []entity.PostReaction, errPostReactions chan error) {
	tempPostReactions, err := u.postReactionsRepo.FetchByUserId(id, like)
	postReactions <- tempPostReactions
	errPostReactions <- err
}

func (u *UsersUsecase) fetchCommentReactions(id int, like bool, commentReactions chan []entity.CommentReaction, errCommentReactions chan error) {
	tempCommentReactions, err := u.commentReactionsRepo.FetchByUserId(id, like)
	commentReactions <- tempCommentReactions
	errCommentReactions <- err
}

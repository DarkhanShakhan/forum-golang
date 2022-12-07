package usecase

import (
	"context"
	"forum_app/internal/entity"
	"log"
)

type UsersUsecase struct {
	userRepo             UsersRepository
	postRepo             PostsRepository
	postReactionsRepo    PostReactionsRepository
	commentRepo          CommentRepository
	commentReactionsRepo CommentReactionsRepository
	errorLog             *log.Logger
}

func NewUsersUsecase(userRepo UsersRepository, postRepo PostsRepository, postReactionsRepo PostReactionsRepository, commentRepo CommentRepository, commentReactionsRepo CommentReactionsRepository, errorLog *log.Logger) *UsersUsecase {
	return &UsersUsecase{
		userRepo:             userRepo,
		postRepo:             postRepo,
		postReactionsRepo:    postReactionsRepo,
		commentRepo:          commentRepo,
		commentReactionsRepo: commentReactionsRepo,
		errorLog:             errorLog,
	}
}

func (u *UsersUsecase) FetchById(ctx context.Context, id int, userRes chan entity.UserResult) {
	user, err := u.userRepo.FetchById(ctx, id)
	if err != nil {
		u.errorLog.Println(err)
		userRes <- entity.UserResult{Err: err}
	}
	u.fetchUserDetails(ctx, &user)
	userRes <- entity.UserResult{User: user}
}

func (u *UsersUsecase) fetchUserDetails(ctx context.Context, user *entity.User) {
	var (
		err                error
		posts              = make(chan []entity.Post)
		comments           = make(chan []entity.Comment)
		postLikes          = make(chan []entity.PostReaction)
		postDislikes       = make(chan []entity.PostReaction)
		commentLikes       = make(chan []entity.CommentReaction)
		commentDislikes    = make(chan []entity.CommentReaction)
		errPosts           = make(chan error)
		errComments        = make(chan error)
		errPostLikes       = make(chan error)
		errPostDislikes    = make(chan error)
		errCommentLikes    = make(chan error)
		errCommentDislikes = make(chan error)
	)
	go u.fetchPosts(ctx, user.Id, posts, errPosts)
	go u.fetchComments(ctx, user.Id, comments, errComments)
	go u.fetchPostReactions(ctx, user.Id, true, postLikes, errPostLikes)
	go u.fetchPostReactions(ctx, user.Id, false, postDislikes, errPostDislikes)
	go u.fetchCommentReactions(ctx, user.Id, true, commentLikes, errCommentLikes)
	go u.fetchCommentReactions(ctx, user.Id, false, commentDislikes, errCommentDislikes)
	for i := 0; i < 6; i++ {
		select {
		case user.Posts = <-posts:
			if err = <-errPosts; err != nil {
				u.errorLog.Println(err)
			}
		case user.Comments = <-comments:
			if err = <-errComments; err != nil {
				u.errorLog.Println(err)
			}
		case user.PostLikes = <-postLikes:
			if err = <-errPostLikes; err != nil {
				u.errorLog.Println(err)
			}
		case user.PostDislikes = <-postDislikes:
			if err = <-errPostDislikes; err != nil {
				u.errorLog.Println(err)
			}
		case user.CommentLikes = <-commentLikes:
			if err = <-errCommentLikes; err != nil {
				u.errorLog.Println(err)
			}

		case user.CommentDislikes = <-commentDislikes:
			if err = <-errCommentDislikes; err != nil {
				u.errorLog.Println(err)
			}
		}
	}
	user.CountTotals()
}

func (u *UsersUsecase) FetchByEmail(ctx context.Context, email string, userRes chan entity.UserResult) {
	user, err := u.userRepo.FetchByEmail(ctx, email)
	if err != nil {
		userRes <- entity.UserResult{Err: err}
	}
	userRes <- entity.UserResult{User: user}
}

func (u *UsersUsecase) FetchAll(ctx context.Context, usersRes chan entity.UsersResult) {
	users, err := u.userRepo.FetchAll(ctx)
	if err != nil {
		usersRes <- entity.UsersResult{Err: err}
	}
	usersRes <- entity.UsersResult{Users: users}
}

// func (u *UsersUsecase) Update(user entity.User) error {
// 	return u.userRepo.Update(user)
// }

// func (u *UsersUsecase) DeleteById(id int) error {
// 	return u.userRepo.DeleteById(id)
// }

func (u *UsersUsecase) fetchPosts(ctx context.Context, id int, posts chan []entity.Post, errPosts chan error) {
	tempPosts, err := u.postRepo.FetchByUserId(ctx, id)
	posts <- tempPosts
	errPosts <- err
}

func (u *UsersUsecase) fetchComments(ctx context.Context, id int, comments chan []entity.Comment, errComments chan error) {
	tempComments, err := u.commentRepo.FetchByUserId(ctx, id)
	comments <- tempComments
	errComments <- err
}

func (u *UsersUsecase) fetchPostReactions(ctx context.Context, id int, like bool, postReactions chan []entity.PostReaction, errPostReactions chan error) {
	tempPostReactions, err := u.postReactionsRepo.FetchByUserId(ctx, id, like)
	postReactions <- tempPostReactions
	errPostReactions <- err
}

func (u *UsersUsecase) fetchCommentReactions(ctx context.Context, id int, like bool, commentReactions chan []entity.CommentReaction, errCommentReactions chan error) {
	tempCommentReactions, err := u.commentReactionsRepo.FetchByUserId(ctx, id, like)
	commentReactions <- tempCommentReactions
	errCommentReactions <- err
}

func (u *UsersUsecase) Store(ctx context.Context, user entity.User, result chan entity.Result) {
	id, err := u.userRepo.Store(ctx, user)
	result <- entity.Result{Id: id, Err: err}
}

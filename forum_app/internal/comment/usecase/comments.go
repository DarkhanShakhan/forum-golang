package usecase

import (
	"context"
	"forum_app/internal/entity"
	"log"
)

type CommentsUsecase struct {
	commentsRepo         CommentsRepository
	commentReactionsRepo CommentReactionsRepository
	postsRepo            PostsRepository
	usersRepo            UsersRepository
	errorLog             *log.Logger
}

func NewCommentsUsecase(commentsRepo CommentsRepository, commentReactionsRepo CommentReactionsRepository, postsRepo PostsRepository, usersRepo UsersRepository, errorLog *log.Logger) *CommentsUsecase {
	return &CommentsUsecase{
		commentsRepo:         commentsRepo,
		commentReactionsRepo: commentReactionsRepo,
		postsRepo:            postsRepo,
		usersRepo:            usersRepo,
		errorLog:             errorLog,
	}
}

func (cu *CommentsUsecase) FetchById(ctx context.Context, id int, commentRes chan entity.CommentResult) {
	comment, err := cu.commentsRepo.FetchById(ctx, id)
	if err != nil {
		commentRes <- entity.CommentResult{Err: err}
	}

	comment.CountTotals()
	commentRes <- entity.CommentResult{Comment: comment}
}

func (cu *CommentsUsecase) fetchCommentDetails(ctx context.Context, comment *entity.Comment) {
	var (
		err         error
		post        = make(chan entity.Post)
		user        = make(chan entity.User)
		likes       = make(chan []entity.Reaction)
		dislikes    = make(chan []entity.Reaction)
		errPost     = make(chan error)
		errUser     = make(chan error)
		errLikes    = make(chan error)
		errDislikes = make(chan error)
	)
	go cu.fetchPost(ctx, comment.Post.Id, post, errPost)
	go cu.fetchUser(ctx, comment.User.Id, user, errUser)
	go cu.fetchReaction(ctx, comment.Id, true, likes, errLikes)
	go cu.fetchReaction(ctx, comment.Id, false, dislikes, errDislikes)

	for i := 0; i < 4; i++ {
		select {
		case comment.Post = <-post:
			if err = <-errPost; err != nil {
				cu.errorLog.Println(err)
			}

		case comment.User = <-user:
			if err = <-errUser; err != nil {
				cu.errorLog.Println(err)
			}
		case comment.Likes = <-likes:
			if err = <-errLikes; err != nil {
				cu.errorLog.Println(err)
			}
		case comment.Dislikes = <-dislikes:
			if err = <-errDislikes; err != nil {
				cu.errorLog.Println(err)
			}
		}
	}
}

func (cu *CommentsUsecase) fetchPost(ctx context.Context, postId int, post chan entity.Post, errPost chan error) {
	tempPost, err := cu.postsRepo.FetchById(ctx, postId)
	post <- tempPost
	errPost <- err
}

func (cu *CommentsUsecase) fetchUser(ctx context.Context, userId int, user chan entity.User, errUser chan error) {
	tempUser, err := cu.usersRepo.FetchById(ctx, userId)
	user <- tempUser
	errUser <- err
}

func (cu *CommentsUsecase) fetchReaction(ctx context.Context, id int, like bool, reactions chan []entity.Reaction, errReactions chan error) {
	tempReactions, err := cu.commentReactionsRepo.FetchByCommentId(ctx, id, like)
	reactions <- tempReactions
	errReactions <- err
}

func (u *CommentsUsecase) FetchReactions(ctx context.Context, id int, reactionsChan chan entity.ReactionsResult) {
	likes, err := u.commentReactionsRepo.FetchByCommentId(ctx, id, true)
	if err != nil {
		reactionsChan <- entity.ReactionsResult{Err: err}
		return
	}
	dislikes, err := u.commentReactionsRepo.FetchByCommentId(ctx, id, false)
	if err != nil {
		reactionsChan <- entity.ReactionsResult{Err: err}
		return
	}
	reactionsChan <- entity.ReactionsResult{Reactions: append(likes, dislikes...)}
}

func (cu *CommentsUsecase) Store(ctx context.Context, comment entity.Comment, res chan entity.Result) {
	id, err := cu.commentsRepo.Store(ctx, comment)
	if err != nil {
		res <- entity.Result{Err: err}
	}
	res <- entity.Result{Id: id}
}

func (cu *CommentsUsecase) StoreCommentReaction(ctx context.Context, commentReaction entity.CommentReaction, err chan error) {
	err <- cu.commentReactionsRepo.StoreReaction(ctx, commentReaction)
}

func (u *CommentsUsecase) UpdateCommentReaction(ctx context.Context, commentReaction entity.CommentReaction, err chan error) {
	err <- u.commentReactionsRepo.UpdateReaction(ctx, commentReaction)
}

func (u *CommentsUsecase) DeleteCommentReaction(ctx context.Context, commentReaction entity.CommentReaction, err chan error) {
	err <- u.commentReactionsRepo.DeleteReaction(ctx, commentReaction)
}

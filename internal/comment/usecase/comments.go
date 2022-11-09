package usecase

import (
	"forum/internal/entity"
	"log"
)

type CommentsUsecase struct {
	commentsRepo         CommentsRepository
	commentReactionsRepo CommentReactionsRepository
	postsRepo            PostsRepository
	usersRepo            UsersRepository
}

func NewCommentsUsecase(commentsRepo CommentsRepository, commentReactionsRepo CommentReactionsRepository, postsRepo PostsRepository, usersRepo UsersRepository) *CommentsUsecase {
	return &CommentsUsecase{commentsRepo: commentsRepo, commentReactionsRepo: commentReactionsRepo, postsRepo: postsRepo, usersRepo: usersRepo}
}

func (cu *CommentsUsecase) FetchById(id int) (entity.Comment, error) {
	comment, err := cu.commentsRepo.FetchById(id)
	if err != nil {
		return entity.Comment{}, err
	}
	post := make(chan entity.Post)
	user := make(chan entity.User)
	likes := make(chan []entity.Reaction)
	dislikes := make(chan []entity.Reaction)
	errPost := make(chan error)
	errUser := make(chan error)
	errLikes := make(chan error)
	errDislikes := make(chan error)
	go cu.fetchPost(id, post, errPost)
	go cu.fetchUser(id, user, errUser)
	go cu.fetchReaction(id, true, likes, errLikes)
	go cu.fetchReaction(id, false, dislikes, errDislikes)

	for i := 0; i < 4; i++ {
		select {
		case comment.Post = <-post:
			if err = <-errPost; err != nil {
				log.Println(err)
			}

		case comment.User = <-user:
			if err = <-errUser; err != nil {
				log.Println(err)
			}
		case comment.Likes = <-likes:
			if err = <-errLikes; err != nil {
				log.Println(err)
			}
		case comment.Dislikes = <-dislikes:
			if err = <-errDislikes; err != nil {
				log.Println(err)
			}
		}
	}
	comment.CountTotals()
	return comment, nil
}

func (cu *CommentsUsecase) fetchPost(id int, post chan entity.Post, errPost chan error) {
	tempPost, err := cu.postsRepo.FetchByCommentId(id)
	post <- tempPost
	errPost <- err
}

func (cu *CommentsUsecase) fetchUser(id int, user chan entity.User, errUser chan error) {
	tempUser, err := cu.usersRepo.FetchByCommentId(id)
	user <- tempUser
	errUser <- err
}

func (cu *CommentsUsecase) fetchReaction(id int, like bool, reactions chan []entity.Reaction, errReactions chan error) {
	tempReactions, err := cu.commentReactionsRepo.FetchByCommentId(id, like)
	reactions <- tempReactions
	errReactions <- err
}

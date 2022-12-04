package usecase

import (
	"forum_app/internal/entity"
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
	go cu.fetchPost(comment.Post.Id, post, errPost)
	go cu.fetchUser(comment.User.Id, user, errUser)
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

func (cu *CommentsUsecase) fetchPost(postId int, post chan entity.Post, errPost chan error) {
	tempPost, err := cu.postsRepo.FetchById(postId)
	post <- tempPost
	errPost <- err
}

func (cu *CommentsUsecase) fetchUser(userId int, user chan entity.User, errUser chan error) {
	tempUser, err := cu.usersRepo.FetchById(userId)
	user <- tempUser
	errUser <- err
}

func (cu *CommentsUsecase) fetchReaction(id int, like bool, reactions chan []entity.Reaction, errReactions chan error) {
	tempReactions, err := cu.commentReactionsRepo.FetchByCommentId(id, like)
	reactions <- tempReactions
	errReactions <- err
}

//create comment
func (cu *CommentsUsecase) Store(comment entity.Comment) (int64, error) {

	id, err := cu.commentsRepo.Store(comment)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// func (cu *CommentsUsecase) Update(comment entity.Comment) error {
// 	return cu.commentsRepo.Update(comment)
// }

// func (cu *CommentsUsecase) DeleteById(id int) error {
// 	return cu.commentsRepo.DeleteById(id)
// }

func (cu *CommentsUsecase) StoreCommentReaction(commentReaction entity.CommentReaction) error {
	return cu.commentReactionsRepo.StoreReaction(commentReaction)
}

func (u *CommentsUsecase) UpdateCommentReaction(commentReaction entity.CommentReaction) error {
	return u.commentReactionsRepo.UpdateReaction(commentReaction)
}

func (u *CommentsUsecase) DeleteCommentReaction(commentReaction entity.CommentReaction) error {
	return u.commentReactionsRepo.DeleteReaction(commentReaction)
}

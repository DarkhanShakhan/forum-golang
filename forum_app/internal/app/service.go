package app

import (
	"context"
	cr "forum_app/internal/comment/repository"
	cUcse "forum_app/internal/comment/usecase"
	pr "forum_app/internal/post/repository"
	pUcse "forum_app/internal/post/usecase"
	ur "forum_app/internal/user/repository"
	uUcse "forum_app/internal/user/usecase"
	"forum_app/pkg/sqlite3"
	"log"
	"time"
)

const duration = 5 * time.Second

type Handler struct {
	errLog  *log.Logger
	infoLog *log.Logger
	ucase   UserUsecase
	pcase   PostUsecase
	ccase   CommentUsecase
}

// FIXME:deal with error from sqlite3
func NewHandler(errLog, infoLog *log.Logger) *Handler {
	db, _ := sqlite3.New()
	usersRepo := ur.NewUsersRepository(db, errLog)
	postsRepo := pr.NewPostsRepository(db, errLog)
	pReactionsRepo := pr.NewPostReactionsRepository(db, errLog)
	categoriesRepo := pr.NewCategoriesRepository(db, errLog)
	commentsRepo := cr.NewCommentsRepository(db, errLog)
	cReactionsRepo := cr.NewCommentReactionsRepository(db, errLog)
	ucase := uUcse.NewUsersUsecase(usersRepo, postsRepo, pReactionsRepo, commentsRepo, cReactionsRepo, errLog)
	pcase := pUcse.NewPostsUsecase(postsRepo, pReactionsRepo, commentsRepo, cReactionsRepo, categoriesRepo, usersRepo, errLog)
	ccase := cUcse.NewCommentsUsecase(commentsRepo, cReactionsRepo, postsRepo, usersRepo, errLog)
	return &Handler{errLog, infoLog, ucase, pcase, ccase}
}

func getTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadline, ok := ctx.Deadline(); ok {
		return context.WithDeadline(context.Background(), deadline)
	}
	return context.WithTimeout(context.Background(), duration)
}

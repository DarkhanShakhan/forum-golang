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
	errorLog *log.Logger
	ucase    UserUsecase
	pcase    PostUsecase
	ccase    CommentUsecase
}

// FIXME:deal with error from sqlite3
func NewHandler(errorLog *log.Logger) *Handler {
	db, _ := sqlite3.New()
	usersRepo := ur.NewUsersRepository(db, errorLog)
	postsRepo := pr.NewPostsRepository(db, errorLog)
	pReactionsRepo := pr.NewPostReactionsRepository(db, errorLog)
	categoriesRepo := pr.NewCategoriesRepository(db, errorLog)
	commentsRepo := cr.NewCommentsRepository(db, errorLog)
	cReactionsRepo := cr.NewCommentReactionsRepository(db, errorLog)
	ucase := uUcse.NewUsersUsecase(usersRepo, postsRepo, pReactionsRepo, commentsRepo, cReactionsRepo, errorLog)
	pcase := pUcse.NewPostsUsecase(postsRepo, pReactionsRepo, commentsRepo, categoriesRepo, usersRepo, errorLog)
	ccase := cUcse.NewCommentsUsecase(commentsRepo, cReactionsRepo, postsRepo, usersRepo, errorLog)
	return &Handler{errorLog, ucase, pcase, ccase}
}

func getTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadline, ok := ctx.Deadline(); ok {
		return context.WithDeadline(context.Background(), deadline)
	}
	return context.WithTimeout(context.Background(), duration)
}

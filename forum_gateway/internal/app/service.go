package app

import (
	"context"
	"forum_gateway/internal/entity"
	"log"
	"time"
)

const duration = 10 * time.Second

type Handler struct {
	errLog     *log.Logger
	infoLog    *log.Logger
	auUcase    AuthUsecase
	forumUcase ForumUsecase
}

func NewHandler(errLog, infoLog *log.Logger, auUcase AuthUsecase, forumUcase ForumUsecase) *Handler {
	return &Handler{errLog, infoLog, auUcase, forumUcase}
}

func getTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadline, ok := ctx.Deadline(); ok {
		return context.WithDeadline(context.Background(), deadline)
	}
	return context.WithTimeout(context.Background(), duration)
}

type AuthUsecase interface {
	SignUp(context.Context, entity.Credentials, chan error)
	SignIn(context.Context, entity.Credentials, chan entity.SessionResult)
	SignOut(context.Context, entity.Session, chan error)
	Authenticate(context.Context, string, chan entity.AuthStatusResult)
}

type ForumUsecase interface {
	FetchPosts(context.Context, chan entity.Response)
	FetchUsers(context.Context, chan entity.Response)
	FetchPost(context.Context, int, chan entity.Response)
	FetchUser(context.Context, int, chan entity.Response)
	FetchCategory(context.Context, int, chan entity.Response)
	StorePost(context.Context, entity.Post, chan entity.Result)
	StoreComment(context.Context, entity.Comment, chan entity.Result)
	PostReaction(context.Context, entity.PostReaction, chan error)
}

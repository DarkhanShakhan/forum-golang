package app

import (
	"context"
	"forum_gateway/internal/entity"
	"log"
	"time"
)

const duration = 10 * time.Second

type Handler struct {
	errLog  *log.Logger
	infoLog *log.Logger
	auUcase AuthUsecase
}

func NewHandler(errLog, infoLog *log.Logger, auUcase AuthUsecase) *Handler {
	return &Handler{errLog, infoLog, auUcase}
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
	Authenticate(context.Context, string, chan entity.AuthStatusResult)
}

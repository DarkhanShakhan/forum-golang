package app

import (
	"context"
	"log"
	"time"
)

const duration = 10 * time.Second

type Handler struct {
	errLog  *log.Logger
	infoLog *log.Logger
}

func NewHandler(errLog, infoLog *log.Logger) *Handler {
	return &Handler{errLog, infoLog}
}

func getTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadline, ok := ctx.Deadline(); ok {
		return context.WithDeadline(context.Background(), deadline)
	}
	return context.WithTimeout(context.Background(), duration)
}

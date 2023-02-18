package app

import (
	"context"
	"forum_gateway/internal/entity"
	"log"
	"time"
)

const duration = 10 * time.Second

type Handler struct {
	errLog      *log.Logger
	infoLog     *log.Logger
	auUcase     AuthUsecase
	forumUcase  ForumUsecase
	config      *Config
	oauths      map[method]OAuth
	rateLimiter *IPRateLimiter
}

func NewHandler(errLog, infoLog *log.Logger, auUcase AuthUsecase, forumUcase ForumUsecase) *Handler {
	h := Handler{
		errLog:      errLog,
		infoLog:     infoLog,
		auUcase:     auUcase,
		forumUcase:  forumUcase,
		config:      NewConfig(),
		oauths:      map[method]OAuth{},
		rateLimiter: NewIPRateLimiter(1, 5),
	}
	h.setOauth([]method{github, google})
	return &h
}

func (h *Handler) setOauth(methods []method) {
	for _, m := range methods {
		clientId, clientSecret := h.getOAuthConfig(m)
		temp, err := NewOAuth(clientId, clientSecret, "pseudo_random", m)
		if err != nil {
			h.errLog.Println(err)
		} else {
			h.oauths[m] = temp
		}
	}
}

func (h *Handler) getOAuthConfig(m method) (string, string) {
	switch m {
	case github:
		return h.config.GitHub.ClientId, h.config.GitHub.ClientSecret
	case google:
		return h.config.Google.ClientId, h.config.Google.ClientSecret
	default:
		return "", ""
	}
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
	OAuth(context.Context, entity.Credentials, chan entity.SessionResult)
}

type ForumUsecase interface {
	FetchPosts(context.Context, chan entity.Response)
	FetchUsers(context.Context, chan entity.Response)
	FetchPost(context.Context, int, chan entity.Response)
	FetchUser(context.Context, int, chan entity.Response)
	FetchCategories(context.Context, chan entity.Response)
	FetchCategory(context.Context, int, chan entity.Response)
	StorePost(context.Context, entity.Post, chan entity.Result)
	StoreComment(context.Context, entity.Comment, chan entity.Result)
	PostReaction(context.Context, entity.PostReaction, chan error)
	CommentReaction(context.Context, entity.CommentReaction, chan error)
}

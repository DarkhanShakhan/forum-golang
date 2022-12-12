package usecase

import (
	"context"
	"forum_auth/internal/entity"
)

type SessionsRepo interface {
	Fetch(ctx context.Context, token string) (entity.Session, error)
	Store(ctx context.Context, session entity.Session) error
	Update(ctx context.Context, session entity.Session) error
	Delete(ctx context.Context, session entity.Session) error
}

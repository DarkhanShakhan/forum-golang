package app

import (
	"context"
	"forum_auth/internal/entity"
)

type AuthUsecase interface {
	SignIn(ctx context.Context, credentials entity.Credentials, sessionRes chan entity.SessionResult)
	SignUp(ctx context.Context, credentials entity.Credentials, credsRes chan entity.CredentialsResult)
	Authenticate(ctx context.Context, session entity.Session, authStatus chan entity.AuthStatusResult)
	SignOut(ctx context.Context, session entity.Session, err chan error)
}

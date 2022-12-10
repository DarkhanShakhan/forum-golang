package auth_usecase

import (
	"context"
	"forum_auth/internal/entity"
	"log"
)

type AuthUsecase struct {
	sessionRepo SessionsRepo
	errLog      *log.Logger
}

func NewAuthUsecase(sessionRepo SessionsRepo, errLog *log.Logger) AuthUsecase {
	return AuthUsecase{sessionRepo: sessionRepo, errLog: errLog}
}

func (au *AuthUsecase) SignIn(ctx context.Context) {}

func (au *AuthUsecase) SignUp() {}

func (au *AuthUsecase) Authenticate() {}

func (au *AuthUsecase) SignOut(ctx context.Context, session entity.Session, err chan error) {
	err <- au.sessionRepo.Delete(ctx, session)
}

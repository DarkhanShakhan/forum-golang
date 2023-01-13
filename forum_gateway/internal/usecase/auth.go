package usecase

import (
	"context"
	"forum_gateway/internal/entity"
	"log"
)

type AuthUsecase struct {
	errLog  *log.Logger
	infoLog *log.Logger
}

func NewAuthUsecase(errLog, infoLog *log.Logger) *AuthUsecase {
	return &AuthUsecase{errLog: errLog, infoLog: infoLog}
}

func (au *AuthUsecase) SignUp(ctx context.Context, credentials entity.Credentials, credsRes chan entity.CredentialsResult) {
	ctx, cancel := get
}

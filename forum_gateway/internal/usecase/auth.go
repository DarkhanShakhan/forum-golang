package usecase

import (
	"context"
	"encoding/json"
	"forum_gateway/internal/entity"
	"log"
	"net/http"
)

type AuthUsecase struct {
	errLog  *log.Logger
	infoLog *log.Logger
}

func NewAuthUsecase(errLog, infoLog *log.Logger) *AuthUsecase {
	return &AuthUsecase{errLog: errLog, infoLog: infoLog}
}

func (au *AuthUsecase) SignUp(ctx context.Context, credentials entity.Credentials, errRes chan error) {
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		errRes <- entity.ErrInternalServer
		return
	}
	response, err := getAPIResponse(ctx, http.MethodPost, "http://localhost:8081/sign_up", requestBody)
	if err != nil {
		errRes <- entity.ErrInternalServer
	}
	switch response.StatusCode {
	case 405, 500:
		errRes <- entity.ErrInternalServer
		return
	case 408:
		errRes <- entity.ErrRequestTimeout
	case 400:
		r, err := getResponse(response.Body)
		if err != nil {
			errRes <- entity.ErrInternalServer
			return
		}
		if r.ErrorMessage == "User with a given email already exists" {
			errRes <- entity.ErrEmailExists
			return
		}
		errRes <- entity.ErrInternalServer
		return
	}
	errRes <- nil
}

func (au *AuthUsecase) SignIn(ctx context.Context, credentials entity.Credentials, sessionChan chan entity.SessionResult) {
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		sessionChan <- entity.SessionResult{Err: err}
		return
	}
	response, err := getAPIResponse(ctx, http.MethodPost, "http://localhost:8081/sign_up", requestBody)
	if err != nil {
		sessionChan <- entity.SessionResult{Err: entity.ErrInternalServer}
	}
	switch response.StatusCode {
	case 400, 500:
		sessionChan <- entity.SessionResult{Err: entity.ErrInternalServer}
	case 404:
		sessionChan <- entity.SessionResult{Err: entity.ErrNotFound}
	case 408:
		sessionChan <- entity.SessionResult{Err: entity.ErrRequestTimeout}
	}
	session, err := getSession(response.Body)
	sessionChan <- entity.SessionResult{Session: session, Err: err}
}

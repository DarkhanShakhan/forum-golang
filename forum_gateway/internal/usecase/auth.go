package usecase

import (
	"context"
	"encoding/json"
	"fmt"
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

func (au *AuthUsecase) SignUp(ctx context.Context, credentials entity.Credentials, errChan chan error) {
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		errChan <- entity.ErrInternalServer
		return
	}
	response, err := getAPIResponse(ctx, http.MethodPost, "http://localhost:8081/sign_up", requestBody)
	if err != nil {
		errChan <- entity.ErrInternalServer
		return
	}
	switch response.StatusCode {
	case 405, 500:
		errChan <- entity.ErrInternalServer
		return
	case 408:
		errChan <- entity.ErrRequestTimeout
		return
	case 400:
		r, err := getResponse(response.Body)
		if err != nil {
			errChan <- entity.ErrInternalServer
			return
		}
		if r.ErrorMessage == "User with a given email already exists" {
			errChan <- entity.ErrEmailExists
			return
		}
		errChan <- entity.ErrInternalServer
		return
	}
	errChan <- nil
}

func (au *AuthUsecase) SignIn(ctx context.Context, credentials entity.Credentials, sessionChan chan entity.SessionResult) {
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		au.errLog.Println(err)
		sessionChan <- entity.SessionResult{Err: err}
		return
	}
	response, err := getAPIResponse(ctx, http.MethodPost, "http://localhost:8081/sign_in", requestBody)
	if err != nil {
		au.errLog.Println(err)
		sessionChan <- entity.SessionResult{Err: entity.ErrInternalServer}
		return
	}
	switch response.StatusCode {
	case 500:
		sessionChan <- entity.SessionResult{Err: entity.ErrInternalServer}
		return
	case 404:
		sessionChan <- entity.SessionResult{Err: entity.ErrNotFound}
		return
	case 408:
		sessionChan <- entity.SessionResult{Err: entity.ErrRequestTimeout}
		return
	case 400:
		r, _ := getResponse(response.Body) // err checking omitted because it returns Internal Server Error anyway
		if r.ErrorMessage == "Invalid Password" {
			sessionChan <- entity.SessionResult{Err: entity.ErrInvalidPassword}
			return
		}
		sessionChan <- entity.SessionResult{Err: entity.ErrInternalServer}
		return
	}
	session, err := getSession(response.Body)
	sessionChan <- entity.SessionResult{Session: session, Err: err}
}

func (au *AuthUsecase) Authenticate(ctx context.Context, token string, authChan chan entity.AuthStatusResult) {
	if token == "" {
		authChan <- entity.AuthStatusResult{Status: entity.NonAuthorised}
		return
	}
	response, err := getAPIResponse(ctx, http.MethodGet, "http://localhost:8081/authenticate", []byte(fmt.Sprintf(`{"token":"%s"}`, token)))
	if err != nil {
		authChan <- entity.AuthStatusResult{Status: entity.NonAuthorised}
		return
	}
	switch response.StatusCode {
	case 408:
		authChan <- entity.AuthStatusResult{Err: entity.ErrRequestTimeout}
	case 200:
		authChan <- getAuthStatus(response.Body)
	default:
		authChan <- entity.AuthStatusResult{Err: entity.ErrInternalServer}
	}
}

func (au *AuthUsecase) SignOut(ctx context.Context, session entity.Session, errChan chan error) {
	requestBody, err := json.Marshal(session)
	if err != nil {
		errChan <- entity.ErrInternalServer
		return
	}
	response, err := getAPIResponse(ctx, http.MethodDelete, "http://localhost:8081/sign_out", requestBody)
	if err != nil {
		errChan <- entity.ErrInternalServer
		return
	}
	switch response.StatusCode {
	case 408:
		errChan <- entity.ErrRequestTimeout
		return
	case 204:
		errChan <- nil
		return
	default: // FIXME: 400??
		errChan <- entity.ErrInternalServer
		return
	}
}

func (au *AuthUsecase) OAuth(ctx context.Context, credentials entity.Credentials, sessionChan chan entity.SessionResult) {
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		au.errLog.Println(err)
		sessionChan <- entity.SessionResult{Err: err}
		return
	}
	response, err := getAPIResponse(ctx, http.MethodPost, "http://localhost:8081/oauth_signin", requestBody)
	if err != nil {
		au.errLog.Println(err)
		sessionChan <- entity.SessionResult{Err: entity.ErrInternalServer}
		return
	}
	switch response.StatusCode {
	case 500:
		sessionChan <- entity.SessionResult{Err: entity.ErrInternalServer}
		return
	case 404:
		sessionChan <- entity.SessionResult{Err: entity.ErrNotFound}
		return
	case 408:
		sessionChan <- entity.SessionResult{Err: entity.ErrRequestTimeout}
		return
	case 400:
		r, _ := getResponse(response.Body) // err checking omitted because it returns Internal Server Error anyway
		if r.ErrorMessage == "Invalid Password" {
			sessionChan <- entity.SessionResult{Err: entity.ErrInvalidPassword}
			return
		}
		sessionChan <- entity.SessionResult{Err: entity.ErrInternalServer}
		return
	}
	session, err := getSession(response.Body)
	sessionChan <- entity.SessionResult{Session: session, Err: err}
}

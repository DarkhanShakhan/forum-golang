package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"forum_auth/internal/entity"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	sessionRepo SessionsRepo
	errLog      *log.Logger
}

func NewAuthUsecase(sessionRepo SessionsRepo, errLog *log.Logger) *AuthUsecase {
	return &AuthUsecase{sessionRepo: sessionRepo, errLog: errLog}
}

func (au *AuthUsecase) SignIn(ctx context.Context, credentials entity.Credentials, sessionRes chan entity.SessionResult) {
	response, err := getAPIResponse(ctx, http.MethodGet, fmt.Sprintf("http://localhost:8080/user/email?email=%s", credentials.Email), nil)
	if err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	switch response.StatusCode {
	case 405, 400, 500:
		sessionRes <- entity.SessionResult{Err: entity.ErrInternalServer}
		return
	case 408:
		sessionRes <- entity.SessionResult{Err: entity.ErrRequestTimeout}
		return
	case 404:
		sessionRes <- entity.SessionResult{Err: entity.ErrNotFound}
		return
	}
	user, err := getUser(response.Body)
	if err != nil {
		sessionRes <- entity.SessionResult{Err: entity.ErrInternalServer}
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		sessionRes <- entity.SessionResult{Err: entity.ErrInvalidPassword}
		return
	}
	token, err := uuid.NewV4()
	if err != nil {
		sessionRes <- entity.SessionResult{Err: entity.ErrInternalServer}
		return
	}
	session := entity.Session{UserId: user.Id, Token: token.String()}
	if session, err = au.sessionRepo.Store(ctx, session); err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	sessionRes <- entity.SessionResult{Session: session}
}

func (au *AuthUsecase) SignUp(ctx context.Context, credentials entity.Credentials, credsRes chan entity.CredentialsResult) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		credsRes <- entity.CredentialsResult{Err: err}
		return
	}
	credentials.Password = string(hashedPassword)
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		credsRes <- entity.CredentialsResult{Err: err}
		return
	}
	response, err := getAPIResponse(ctx, http.MethodPost, "http://localhost:8080/user/save", requestBody)
	if err != nil {
		credsRes <- entity.CredentialsResult{Err: err}
		return
	}
	switch response.StatusCode {
	case 408:
		credsRes <- entity.CredentialsResult{Err: entity.ErrRequestTimeout}
		return
	case 405, 500:
		credsRes <- entity.CredentialsResult{Err: entity.ErrInternalServer}
		return
	case 400:
		r, err := getResponse(response.Body)
		if err != nil {
			credsRes <- entity.CredentialsResult{Err: entity.ErrInternalServer}
			return
		}
		if r.ErrorMessage == "User with a given email already exists" {
			credsRes <- entity.CredentialsResult{Err: entity.ErrEmailExists}
			return
		}
		credsRes <- entity.CredentialsResult{Err: entity.ErrInternalServer}
		return
	}
	user, err := getUser(response.Body)
	if err != nil {
		credsRes <- entity.CredentialsResult{Err: err}
		return
	}
	credsRes <- entity.CredentialsResult{Credentials: user}
}

func (au *AuthUsecase) Authenticate(ctx context.Context, session entity.Session, authStatus chan entity.AuthStatusResult) {
	session, err := au.sessionRepo.Fetch(ctx, session.Token)
	if err != nil {
		authStatus <- entity.AuthStatusResult{Status: entity.NonAuthorised, Err: err}
		return
	}
	if session.Token == "" {
		authStatus <- entity.AuthStatusResult{Status: entity.NonAuthorised, Err: errors.New("session doesn't exist")}
		return
	}
	if time.Now().After(session.ExpiryTime) {
		if err = au.sessionRepo.Delete(ctx, session); err != nil {
			au.errLog.Println(err)
		}
		authStatus <- entity.AuthStatusResult{Status: entity.NonAuthorised, Err: errors.New("session expired")}
		return
	}
	session, err = au.sessionRepo.Update(ctx, session)
	if err != nil {
		authStatus <- entity.AuthStatusResult{Status: entity.NonAuthorised, Err: err}
	}
	authStatus <- entity.AuthStatusResult{Status: entity.Authorised, Session: session}
}

func (au *AuthUsecase) SignOut(ctx context.Context, session entity.Session, err chan error) {
	err <- au.sessionRepo.Delete(ctx, session)
}

func (au *AuthUsecase) OauthSignIn(ctx context.Context, credentials entity.Credentials, sessionRes chan entity.SessionResult) {
	response, err := getAPIResponse(ctx, http.MethodGet, fmt.Sprintf("http://localhost:8080/user/email?email=%s", credentials.Email), nil)
	if err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	user, err := getUser(response.Body)
	if err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	if user.Id == 0 { // if user doesn't exist, it should be stored
		res := storeUser(ctx, credentials)
		if res.Err != nil {
			sessionRes <- entity.SessionResult{Err: res.Err}
			return
		}
		user = res.Credentials
	}
	sessionRes <- au.createSession(ctx, user)
}

func getAPIResponse(ctx context.Context, method string, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	return client.Do(req)
}

func storeUser(ctx context.Context, credentials entity.Credentials) entity.CredentialsResult {
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	response, err := getAPIResponse(ctx, http.MethodPost, "http://localhost:8080/user/save", requestBody)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	user, err := getUser(response.Body)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	return entity.CredentialsResult{Credentials: user}
}

func getUser(response io.ReadCloser) (entity.Credentials, error) {
	temp := entity.Response{}
	err := json.NewDecoder(response).Decode(&temp)
	if err != nil {
		return entity.Credentials{}, err
	}
	jsonUser, err := json.Marshal(temp.Body)
	if err != nil {
		return entity.Credentials{}, err
	}
	user := entity.Credentials{}
	err = json.Unmarshal(jsonUser, &user)
	return user, err
}

func getResponse(response io.ReadCloser) (entity.Response, error) {
	result := entity.Response{}
	err := json.NewDecoder(response).Decode(&result)
	return result, err
}

func (au *AuthUsecase) createSession(ctx context.Context, credentials entity.Credentials) entity.SessionResult {
	token, err := uuid.NewV4()
	if err != nil {
		return entity.SessionResult{Err: err}
	}
	session := entity.Session{UserId: credentials.Id, Token: token.String()}
	if session, err = au.sessionRepo.Store(ctx, session); err != nil {
		return entity.SessionResult{Err: err}
	}
	return entity.SessionResult{Session: session}
}

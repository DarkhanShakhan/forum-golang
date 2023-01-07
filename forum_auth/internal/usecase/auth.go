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
	requestUrl := fmt.Sprintf("http://localhost:8080/user/email?email=%s", credentials.Email)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	client := http.Client{}
	response, err := client.Do(req)
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
		fmt.Println(err)
		fmt.Println(user.Password)
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
		// if strings.Contains(err.Error(), "UNIQUE constraint failed: sessions.user_id") {
		// 	checkExpiryTime()
		// }
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	fmt.Println(session.ExpiryTime)
	sessionRes <- entity.SessionResult{Session: session}
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

func (au *AuthUsecase) SignUp(ctx context.Context, credentials entity.Credentials, credsRes chan entity.CredentialsResult) {
	requestUrl := "http://localhost:8080/user/save"
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl, bytes.NewReader(requestBody))
	if err != nil {
		credsRes <- entity.CredentialsResult{Err: err}
		return
	}
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		credsRes <- entity.CredentialsResult{Err: err}
		return
	}
	user := entity.Credentials{}
	err = json.NewDecoder(response.Body).Decode(&user)
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
	// FIXME:validate expiry date

	token, err := uuid.NewV4()
	if err != nil {
		authStatus <- entity.AuthStatusResult{Status: entity.NonAuthorised, Err: err}
		return
	}
	session.Token = token.String() // updates token: expiry is updated in repo
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
	requestUrl := fmt.Sprintf("http://localhost:8080/user/email?email=%s", credentials.Email)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	user := entity.Credentials{}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	if user.Id == 0 {
		res := au.storeUser(ctx, credentials)
		if res.Err != nil {
			sessionRes <- entity.SessionResult{Err: res.Err}
			return
		}
		user = res.Credentials
	}
	sessionRes <- au.createSession(ctx, user)
}

func (au *AuthUsecase) storeUser(ctx context.Context, credentials entity.Credentials) entity.CredentialsResult {
	requestUrl := "http://localhost:8080/user/save"
	requestBody, err := json.Marshal(credentials)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl, bytes.NewReader(requestBody))
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	user := entity.Credentials{}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return entity.CredentialsResult{Err: err}
	}
	return entity.CredentialsResult{Credentials: user}
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

package auth_usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"forum_auth/internal/entity"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	sessionRepo SessionsRepo
	errLog      *log.Logger
}

func NewAuthUsecase(sessionRepo SessionsRepo, errLog *log.Logger) AuthUsecase {
	return AuthUsecase{sessionRepo: sessionRepo, errLog: errLog}
}

func (au *AuthUsecase) SignIn(ctx context.Context, credentials entity.Credentials, sessionRes chan entity.SessionResult) {
	requestUrl := fmt.Sprintf("localhost:8080/user/email?email=%s", credentials.Email)
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

	//FIXME:validate user
	user := entity.Credentials{}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	token, err := uuid.NewV4()
	if err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	session := entity.Session{UserId: user.Id, Token: token.String()}
	if err = au.sessionRepo.Store(ctx, session); err != nil {
		sessionRes <- entity.SessionResult{Err: err}
		return
	}
	sessionRes <- entity.SessionResult{Session: session}
}

func (au *AuthUsecase) SignUp(ctx context.Context, credentials entity.Credentials, credsRes chan entity.CredentialsResult) {
	requestUrl := "localhost:8080/user/save"
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
	//FIXME:validate expiry date

	token, err := uuid.NewV4()
	if err != nil {
		authStatus <- entity.AuthStatusResult{Status: entity.NonAuthorised, Err: err}
		return
	}
	session.Token = token.String() //updates token: expiry is updated in repo
	err = au.sessionRepo.Update(ctx, session)
	if err != nil {
		authStatus <- entity.AuthStatusResult{Status: entity.NonAuthorised, Err: err}
	}
	authStatus <- entity.AuthStatusResult{Status: entity.Authorised}
}

func (au *AuthUsecase) SignOut(ctx context.Context, session entity.Session, err chan error) {
	err <- au.sessionRepo.Delete(ctx, session)
}

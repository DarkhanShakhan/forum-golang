package app

import (
	"context"
	"encoding/json"
	"fmt"
	"forum_auth/internal/entity"
	"forum_auth/internal/repository"
	"forum_auth/internal/usecase"
	"forum_auth/pkg/sqlite3"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	errorLog *log.Logger
	aucase   AuthUsecase
}

const duration = 5 * time.Second

// FIXME:deal with error from sqlite3
func NewHandler(errorLog *log.Logger) *Handler {
	db, _ := sqlite3.New()
	authRepo := repository.NewSessionsRepository(db, errorLog)
	aucase := usecase.NewAuthUsecase(authRepo, errorLog)
	return &Handler{errorLog: errorLog, aucase: aucase}
}

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if deadline, ok := r.Context().Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), duration)
	}
	defer cancel()
	if r.Method != http.MethodPost {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		w.WriteHeader(405)
		return
	}
	sessionChan := make(chan entity.SessionResult)
	var sessionRes entity.SessionResult
	var err error
	credentials := entity.Credentials{}
	err = json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	go h.aucase.SignIn(ctx, credentials, sessionChan)
	select {
	case sessionRes = <-sessionChan:
		if sessionRes.Err != nil {
			h.errorLog.Println(sessionRes.Err)
			w.WriteHeader(500) // FIXME: no sure
			return
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("{\"token\":%s}", sessionRes.Session.Token)))
}

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if deadline, ok := r.Context().Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), duration)
	}
	defer cancel()
	if r.Method != http.MethodPost {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		w.WriteHeader(405)
		return
	}
	credentialsChan := make(chan entity.CredentialsResult)
	var credentialsRes entity.CredentialsResult
	var err error
	credentials := entity.Credentials{}
	err = json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	go h.aucase.SignUp(ctx, credentials, credentialsChan)
	select {
	case credentialsRes = <-credentialsChan:
		if credentialsRes.Err != nil {
			h.errorLog.Println(credentialsRes.Err)
			w.WriteHeader(500) // FIXME: no sure
			return
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("{\"id\":%d}", credentials.Id))) // is it necessary?
}
func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if deadline, ok := r.Context().Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), duration)
	}
	defer cancel()
	if r.Method != http.MethodDelete {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		w.WriteHeader(405)
		return
	}
	errChan := make(chan error)
	session := entity.Session{}
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	go h.aucase.SignOut(ctx, session, errChan)
	select {
	case err = <-errChan:
		if err != nil {
			h.errorLog.Println(err)
			w.WriteHeader(500) // FIXME: no sure
			return
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204) // no content
}

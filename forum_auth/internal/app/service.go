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
	"strings"
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
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodPost {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	sessionChan := make(chan entity.SessionResult)
	var sessionRes entity.SessionResult
	var err error
	credentials := entity.Credentials{}
	err = json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil || !validateCredentials(credentials) {
		h.errorLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	go h.aucase.SignIn(ctx, credentials, sessionChan)
	select {
	case sessionRes = <-sessionChan:
		err = sessionRes.Err
		if err != nil {
			h.errorLog.Println(err)
			if isConstraintError(err) {
				h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"})
				return
			}
			switch err {
			case entity.ErrNotFound:
				h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "User with a given email doesn't exist"})
			case entity.ErrRequestTimeout:
				h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
			case entity.ErrInvalidPassword:
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Invalid Password"})
			default:
				h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			}
			return
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	}
	h.APIResponse(w, http.StatusCreated, entity.Response{Body: sessionRes.Session})
}

func validateCredentials(credentials entity.Credentials) bool {
	return credentials.Email != "" && credentials.Password != ""
}
func isConstraintError(err error) bool {
	return strings.Contains(err.Error(), "constraint failed")
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
	result, err := json.Marshal(credentials)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // FIXME: not sure
		return
	}
	w.Write(result)
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
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
	if r.Method != http.MethodGet {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		w.WriteHeader(405)
		return
	}
	authStatusChan := make(chan entity.AuthStatusResult)
	authStatusRes := entity.AuthStatusResult{}
	var session entity.Session
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	go h.aucase.Authenticate(ctx, session, authStatusChan)
	select {
	case authStatusRes = <-authStatusChan:
		if authStatusRes.Err != nil {
			h.errorLog.Println(authStatusRes.Err)
			w.Header().Set("Content-Type", "application/json")
			authStatusRes.Err = nil
			result, err := json.Marshal(authStatusRes)
			if err != nil {
				h.errorLog.Println(err)
				w.WriteHeader(500) // FIXME: not sure
				return
			}
			w.WriteHeader(200) // FIXME: not sure
			w.Write(result)
			return
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	}
	w.Header().Set("Content-Type", "application/json")
	result, err := json.Marshal(authStatusRes)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // FIXME: not sure
		return
	}
	w.WriteHeader(200)
	w.Write(result)
}

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

func (h *Handler) OauthSignInHandler(w http.ResponseWriter, r *http.Request) {
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
	go h.aucase.OauthSignIn(ctx, credentials, sessionChan)
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
	result, err := json.Marshal(sessionRes.Session)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // FIXME: not sure
		return
	}
	w.Write(result)
}

func (h *Handler) APIResponse(w http.ResponseWriter, code int, response entity.Response) {
	if code == 204 {
		w.WriteHeader(204)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func getTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadline, ok := ctx.Deadline(); ok {
		return context.WithDeadline(context.Background(), deadline)
	}
	return context.WithTimeout(context.Background(), duration)
}

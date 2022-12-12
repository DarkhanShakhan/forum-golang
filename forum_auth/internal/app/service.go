package app

import (
	"log"
	"net/http"

	"forum_auth/internal/repository"
	"forum_auth/internal/usecase"
	"forum_auth/pkg/sqlite3"
)

type Handler struct {
	errorLog *log.Logger
	aucase   AuthUsecase
}

//FIXME:deal with error from sqlite3
func NewHandler(errorLog *log.Logger) *Handler {
	db, _ := sqlite3.New()
	authRepo := repository.NewSessionsRepository(db, errorLog)
	aucase := usecase.NewAuthUsecase(authRepo, errorLog)
	return &Handler{errorLog: errorLog, aucase: aucase}
}

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request)   {}
func (h *Handler) SignOutHandler(w http.ResponseWriter, r *http.Request) {}

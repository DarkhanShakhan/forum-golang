package app

import (
	"encoding/json"
	"fmt"
	cr "forum_app/internal/comment/repository"
	cUcse "forum_app/internal/comment/usecase"
	"forum_app/internal/entity"
	pr "forum_app/internal/post/repository"
	pUcse "forum_app/internal/post/usecase"
	ur "forum_app/internal/user/repository"
	uUcse "forum_app/internal/user/usecase"
	"forum_app/pkg/sqlite3"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	errorLog *log.Logger
	ucase    UserUsecase
	pcase    PostUsecase
	ccase    CommentUsecase
}

func NewHandler(errorLog *log.Logger) *Handler {
	db, _ := sqlite3.New()
	usersRepo := ur.NewUsersRepository(db)
	postsRepo := pr.NewPostsRepository(db)
	pReactionsRepo := pr.NewPostReactionsRepository(db)
	categoriesRepo := pr.NewCategoriesRepository(db)
	commentsRepo := cr.NewCommentsRepository(db)
	cReactionsRepo := cr.NewCommentReactionsRepository(db)
	ucase := uUcse.NewUsersUsecase(usersRepo, postsRepo, pReactionsRepo, commentsRepo, cReactionsRepo)
	pcase := pUcse.NewPostsUsecase(postsRepo, pReactionsRepo, commentsRepo, categoriesRepo, usersRepo)
	ccase := cUcse.NewCommentsUsecase(commentsRepo, cReactionsRepo, postsRepo, usersRepo)
	return &Handler{errorLog, ucase, pcase, ccase}
}

func (h *Handler) UsersAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorLog.Println("method not allowed")
		w.WriteHeader(405)
		return
	}
	users, err := h.ucase.FetchAll()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		return
	}
	response, err := json.Marshal(users)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}

//FIXME: errorLog
func (h *Handler) UserByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("UsersByIdHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseForm(); err != nil {
		log.Println("UserByIdHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		log.Println("UserByIdHandler: " + err.Error())
		w.WriteHeader(400)
		return
	}
	user, err := h.ucase.FetchById(id)
	if err != nil {
		log.Println("UserByIdHandler: " + err.Error())
		w.WriteHeader(500) //not sure
		return
	}
	response, err := json.Marshal(user)
	if err != nil {
		log.Println("UserByIdHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}

func (h *Handler) UserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("UsersInfoHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseForm(); err != nil {
		log.Println("UserInfoHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}

	email := r.Form.Get("email")
	user, err := h.ucase.FetchByEmail(email)
	if err != nil {
		log.Println("UserByIdHandler: " + err.Error())
		w.WriteHeader(500) //not sure
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		log.Println("UserByIdHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}
func (h *Handler) PostFullHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("PostByIdHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseForm(); err != nil {
		log.Println("PostByIdHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		log.Println("PostByIdHandler: " + err.Error())
		w.WriteHeader(400)
		return
	}
	post, err := h.pcase.FetchById(id)
	response, err := json.Marshal(post)
	if err != nil {
		log.Println("PostByIdHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}

func (h *Handler) PostsAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("PostsAllHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	posts, err := h.pcase.FetchAll()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Println("PostsAllHandler" + err.Error())
		w.WriteHeader(500)
		return
	}
	response, err := json.Marshal(posts)
	if err != nil {
		log.Println("PostsAllHandler" + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}

func (h *Handler) CategoryPostsHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) StorePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Println("StorePostHandler: method not allowed (" + r.Method + ")")
		w.WriteHeader(405)
		return
	}
	var post entity.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Println("StorePostHandler: bad request")
		w.WriteHeader(400)
		return
	}
	id, err := h.pcase.Store(post)
	if err != nil {
		log.Println("StorePostHandler: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("{\"id\":%d}", id)))
}

func (h *Handler) StorePostReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) UpdatePostReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeletePostReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) CommentByIdHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) StoreCommentHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) StoreCommentReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) UpdateCommentReactionHandler(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeleteCommentReaction(w http.ResponseWriter, r *http.Request) {}

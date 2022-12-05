package app

import (
	"context"
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
	"time"
)

const duration = 5 * time.Second

type Handler struct {
	errorLog *log.Logger
	ucase    UserUsecase
	pcase    PostUsecase
	ccase    CommentUsecase
}

// TODO: done
func NewHandler(errorLog *log.Logger) *Handler {
	db, _ := sqlite3.New()
	usersRepo := ur.NewUsersRepository(db, errorLog)
	postsRepo := pr.NewPostsRepository(db, errorLog)
	pReactionsRepo := pr.NewPostReactionsRepository(db, errorLog)
	categoriesRepo := pr.NewCategoriesRepository(db, errorLog)
	commentsRepo := cr.NewCommentsRepository(db, errorLog)
	cReactionsRepo := cr.NewCommentReactionsRepository(db, errorLog)
	ucase := uUcse.NewUsersUsecase(usersRepo, postsRepo, pReactionsRepo, commentsRepo, cReactionsRepo, errorLog)
	pcase := pUcse.NewPostsUsecase(postsRepo, pReactionsRepo, commentsRepo, categoriesRepo, usersRepo, errorLog)
	ccase := cUcse.NewCommentsUsecase(commentsRepo, cReactionsRepo, postsRepo, usersRepo, errorLog)
	return &Handler{errorLog, ucase, pcase, ccase}
}

func (h *Handler) UsersAllHandler(w http.ResponseWriter, r *http.Request) {
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
	users, err := h.ucase.FetchAll(ctx)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	default:
	}
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	response, err := json.Marshal(users)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// func (h *Handler) UserByIdHandler(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		ctx    context.Context
// 		cancel context.CancelFunc
// 	)
// 	if deadline, ok := r.Context().Deadline(); ok {
// 		ctx, cancel = context.WithDeadline(context.Background(), deadline)
// 	} else {
// 		ctx, cancel = context.WithTimeout(context.Background(), duration)
// 	}
// 	defer cancel()
// 	if r.Method != http.MethodGet {
// 		log.Println("UsersByIdHandler: method not allowed (" + r.Method + ")")
// 		w.WriteHeader(405)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := r.ParseForm(); err != nil {
// 		log.Println("UserByIdHandler: " + err.Error())
// 		w.WriteHeader(500)
// 		return
// 	}
// 	id, err := strconv.Atoi(r.Form.Get("id"))
// 	if err != nil {
// 		log.Println("UserByIdHandler: " + err.Error())
// 		w.WriteHeader(400)
// 		return
// 	}
// 	user, err := h.ucase.FetchById(ctx, id)
// 	if err != nil {
// 		log.Println("UserByIdHandler: " + err.Error())
// 		w.WriteHeader(500) //not sure
// 		return
// 	}
// 	response, err := json.Marshal(user)
// 	if err != nil {
// 		log.Println("UserByIdHandler: " + err.Error())
// 		w.WriteHeader(500)
// 		return
// 	}
// 	w.Write(response)
// }

// TODO: done
func (h *Handler) UserByEmailHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(405) // method not allowed
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error ???
		return
	}

	email := r.Form.Get("email")
	if email == "" {
		h.errorLog.Println("bad request: email is not provided")
		w.WriteHeader(400) // bad request
		return
	}
	user, err := h.ucase.FetchByEmail(ctx, email)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	default:
	}
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error ???
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// TODO:done

func (h *Handler) PostDetailsHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := r.ParseForm(); err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(400)
		return
	}
	post, err := h.pcase.FetchById(ctx, id)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	default:
	}
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error ???
		return
	}
	response, err := json.Marshal(post)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// TODO: done
func (h *Handler) PostsAllHandler(w http.ResponseWriter, r *http.Request) {
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
	posts, err := h.pcase.FetchAll(ctx)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	default:
	}
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		return
	}
	response, err := json.Marshal(posts)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// TODO:done

func (h *Handler) UserDetailsHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(405) // method not allowed
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error ???
		return
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(400)
		return
	}
	user, err := h.ucase.FetchById(ctx, id)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	default:
	}
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error ???
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// TODO:
func (h *Handler) CategoryPostsHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(405) // method not allowed
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error ???
		return
	}

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(400)
		return
	}
	category, err := h.pcase.FetchCategoryPosts(ctx, id)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	default:
	}
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error ???
		return
	}

	response, err := json.Marshal(category)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) StorePostHandler(w http.ResponseWriter, r *http.Request) {
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
	id, err := h.pcase.Store(ctx, post)
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

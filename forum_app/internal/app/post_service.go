package app

import (
	"context"
	"encoding/json"
	"fmt"
	"forum_app/internal/entity"
	"log"
	"net/http"
	"strconv"
)

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
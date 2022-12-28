package app

import (
	"context"
	"encoding/json"
	"fmt"
	"forum_app/internal/entity"
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
	postChan := make(chan entity.PostResult)
	var postResult entity.PostResult
	go h.pcase.FetchById(ctx, id, postChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	case postResult = <-postChan:
		if postResult.Err != nil {
			h.errorLog.Println(postResult.Err)
			w.WriteHeader(500) // internal server error ???
			return
		}
	}
	response, err := json.Marshal(entity.Response{Content: postResult.Post})
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

//FIXME: change contexts to goroutines
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
	postsChan := make(chan entity.PostsResult)
	var postsRes entity.PostsResult
	var err error
	go h.pcase.FetchAll(ctx, postsChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	case postsRes = <-postsChan:
		if err = postsRes.Err; err != nil {
			h.errorLog.Println(err)
			w.WriteHeader(500)
			return
		}
	}

	response, err := json.Marshal(postsRes.Posts)
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
	catChan := make(chan entity.CatResult)
	var catResult entity.CatResult
	go h.pcase.FetchCategoryPosts(ctx, id, catChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	case catResult = <-catChan:
		if catResult.Err != nil {
			h.errorLog.Println(catResult.Err)
			w.WriteHeader(500) // internal server error ???
			return
		}
	}

	response, err := json.Marshal(catResult.Cat)
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
	if r.Method != http.MethodPost {
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		w.WriteHeader(405)
		return
	}
	var post entity.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil || !validatePostData(post) {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	resChan := make(chan entity.Result)
	var res entity.Result
	go h.pcase.Store(ctx, post, resChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	case res = <-resChan:
		if res.Err != nil {
			h.errorLog.Println(res.Err)
			w.WriteHeader(500)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("{\"id\":%d}", res.Id)))
}
func validatePostData(post entity.Post) bool {
	if post.Title == "" {
		return false
	} else if post.User.Id == 0 {
		return false
	} else if post.Category == nil {
		return false
	}
	return true
}
func (h *Handler) StorePostReactionHandler(w http.ResponseWriter, r *http.Request) {
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
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		w.WriteHeader(405)
		return
	}
	var post_reaction entity.PostReaction
	err := json.NewDecoder(r.Body).Decode(&post_reaction)
	//FIXME:validate data
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	errChan := make(chan error)
	go h.pcase.StorePostReaction(ctx, post_reaction, errChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	case err = <-errChan:
		if err != nil {
			h.errorLog.Println(err)
			w.WriteHeader(500)
			return
		}
	}
	w.WriteHeader(201)
}
func (h *Handler) UpdatePostReactionHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if _, ok := r.Context().Deadline(); ok {
		ctx, cancel = context.WithCancel(r.Context()) //deadline is already set
	} else {
		ctx, cancel = context.WithTimeout(r.Context(), duration)
	}
	defer cancel()
	if r.Method != http.MethodPut {
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		w.WriteHeader(405)
		return
	}
	var post_reaction entity.PostReaction
	err := json.NewDecoder(r.Body).Decode(&post_reaction)
	//FIXME:validate data
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	errChan := make(chan error)
	go h.pcase.UpdatePostReaction(ctx, post_reaction, errChan)
	select {
	case err = <-errChan:
		if err != nil {
			h.errorLog.Println(err)
			w.WriteHeader(500)
			return
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	}
	w.WriteHeader(204)
}

func (h *Handler) DeletePostReactionHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if _, ok := r.Context().Deadline(); ok {
		ctx, cancel = context.WithCancel(r.Context()) //deadline is already set
	} else {
		ctx, cancel = context.WithTimeout(r.Context(), duration)
	}
	defer cancel()
	if r.Method != http.MethodDelete {
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		w.WriteHeader(405)
		return
	}
	var post_reaction entity.PostReaction
	err := json.NewDecoder(r.Body).Decode(&post_reaction)
	//FIXME:validate data
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	errChan := make(chan error)
	go h.pcase.DeletePostReaction(ctx, post_reaction, errChan)
	select {
	case err = <-errChan:
		if err != nil {
			h.errorLog.Println(err)
			w.WriteHeader(500)
			return
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	}
	w.WriteHeader(204)
}

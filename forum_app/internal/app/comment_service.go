package app

import (
	"context"
	"encoding/json"
	"fmt"
	"forum_app/internal/entity"
	"net/http"
)

func (h *Handler) StoreCommentHandler(w http.ResponseWriter, r *http.Request) {
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
	var comment entity.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	//FIXME:validate data
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	id, err := h.ccase.Store(ctx, comment)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("{\"id\":%d}", id)))
}

func (h *Handler) StoreCommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if deadline, ok := r.Context().Deadline(); ok {
		ctx, cancel = context.WithDeadline(r.Context(), deadline)
	} else {
		ctx, cancel = context.WithTimeout(r.Context(), duration)
	}
	defer cancel()
	if r.Method != http.MethodPost {
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		w.WriteHeader(405)
		return
	}
	var comment_reaction entity.CommentReaction
	err := json.NewDecoder(r.Body).Decode(&comment_reaction)
	//FIXME:validate data
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	err = h.ccase.StoreCommentReaction(ctx, comment_reaction)
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
	w.WriteHeader(201)
}
func (h *Handler) UpdateCommentReactionHandler(w http.ResponseWriter, r *http.Request) {

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if deadline, ok := r.Context().Deadline(); ok {
		ctx, cancel = context.WithDeadline(r.Context(), deadline)
	} else {
		ctx, cancel = context.WithTimeout(r.Context(), duration)
	}
	defer cancel()
	if r.Method != http.MethodPut {
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		w.WriteHeader(405)
		return
	}
	var comment_reaction entity.CommentReaction
	err := json.NewDecoder(r.Body).Decode(&comment_reaction)
	//FIXME:validate data
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	errChan := make(chan error)
	go h.ccase.UpdateCommentReaction(ctx, comment_reaction, errChan)
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

func (h *Handler) DeleteCommentReactionHandler(w http.ResponseWriter, r *http.Request) {
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
	var comment_reaction entity.CommentReaction
	err := json.NewDecoder(r.Body).Decode(&comment_reaction)
	//FIXME:validate data
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	errChan := make(chan error)
	go h.ccase.DeleteCommentReaction(ctx, comment_reaction, errChan)
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

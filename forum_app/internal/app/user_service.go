package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

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
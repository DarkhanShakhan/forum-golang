package app

import (
	"context"
	"encoding/json"
	"fmt"
	"forum_app/internal/entity"
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

	usersChan := make(chan entity.UsersResult)
	var usersRes entity.UsersResult
	var err error
	h.ucase.FetchAll(ctx, usersChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	case usersRes = <-usersChan:

		if err = usersRes.Err; err != nil {
			h.errorLog.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}

	response, err := json.Marshal(usersRes.Users)
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
	userChan := make(chan entity.UserResult)
	var userRes entity.UserResult
	var err error
	go h.ucase.FetchByEmail(ctx, email, userChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	case userRes = <-userChan:
		err = userRes.Err
		if err != nil {
			h.errorLog.Println(err)
			w.WriteHeader(500) // internal server error ???
			return
		}
	}

	response, err := json.Marshal(userRes.User)
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
	userChan := make(chan entity.UserResult)
	var userRes entity.UserResult
	go h.ucase.FetchById(ctx, id, userChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		w.WriteHeader(408) // request timeout
		return
	case userRes = <-userChan:
		err = userRes.Err
		if err != nil {
			h.errorLog.Println(err)
			w.WriteHeader(500) // internal server error ???
			return
		}
	}

	response, err := json.Marshal(userRes.User)
	if err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) StoreUserHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(405) // method not allowed
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorLog.Println(err)
		w.WriteHeader(500) // internal server error ???
		return
	}
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.errorLog.Println("bad request")
		w.WriteHeader(400)
		return
	}
	resChan := make(chan entity.Result)
	var res entity.Result
	go h.ucase.Store(ctx, user, resChan)
	select {
	case res = <-resChan:
		if res.Err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("{\"id\":%d}", res.Id)))
}

package app

import (
	"encoding/json"
	"fmt"
	"forum_app/internal/entity"
	"net/http"
	"strconv"
)

func (h *Handler) UsersAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()

	if r.Method != http.MethodGet {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	usersChan := make(chan entity.UsersResult)
	var usersRes entity.UsersResult
	var err error
	go h.ucase.FetchAll(ctx, usersChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case usersRes = <-usersChan:
		if err = usersRes.Err; err != nil {
			h.errorLog.Println(err)
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}

	h.APIResponse(w, http.StatusOK, entity.Response{Body: usersRes.Users})
}

func (h *Handler) UserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodGet {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}

	email := r.Form.Get("email")
	if email == "" {
		h.errorLog.Println("bad request: email is not provided")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
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
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case userRes = <-userChan:
		err = userRes.Err
		if err != nil {
			h.errorLog.Println(err)
			if err == entity.ErrUserNotFound {
				h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusOK, entity.Response{Body: userRes.User})
}

func (h *Handler) UserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodGet {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	if err := r.ParseForm(); err != nil {
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	userChan := make(chan entity.UserResult)
	var userRes entity.UserResult
	go h.ucase.FetchById(ctx, id, userChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case userRes = <-userChan:
		err = userRes.Err
		if err != nil {
			h.errorLog.Println(err)
			if err == entity.ErrUserNotFound {
				h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusOK, entity.Response{Body: userRes.User})
}

func (h *Handler) StoreUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodPost {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	if err := r.ParseForm(); err != nil {
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.errorLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	resChan := make(chan entity.Result)
	var res entity.Result
	go h.ucase.Store(ctx, user, resChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case res = <-resChan:
		err = res.Err
		if err != nil {
			h.errorLog.Println(err)
			if err == entity.ErrUserExists {
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "User with a given email already exists"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusCreated, entity.Response{Body: entity.User{Id: int(res.Id)}})
}

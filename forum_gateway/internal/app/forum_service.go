package app

import (
	"errors"
	"fmt"
	"forum_gateway/internal/entity"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "templates/errors.html")
		return
	}
	if r.URL.Path != "/" && r.URL.Path != "/posts" {
		h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"}, "templates/errors.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	response := entity.Response{}
	responseChan := make(chan entity.Response)
	go h.forumUcase.FetchPosts(ctx, responseChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case response = <-responseChan:
		err := response.Err
		switch err {
		case entity.ErrInternalServer:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		case nil:
			var (
				auth interface{} = r.Context().Value("authorised")
				id   interface{} = r.Context().Value("user_id")
			)
			response.AuthStatus, _ = auth.(bool)
			user_id, ok := id.(int64)
			if ok {
				response.UserId = user_id
			}
			h.APIResponse(w, http.StatusOK, response, "templates/index.html")
		}
	}
}

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "templates/errors.html")
		return
	}
	post_id, err := getID(r.URL.String(), "posts")
	if err != nil {
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: err.Error()}, "templates/errors.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	response := entity.Response{}
	responseChan := make(chan entity.Response)
	go h.forumUcase.FetchPost(ctx, post_id, responseChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case response = <-responseChan:
		err := response.Err
		switch err {
		case entity.ErrInternalServer:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		case entity.ErrNotFound:
			h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"}, "templates/errors.html")
		case nil:
			var (
				auth interface{} = r.Context().Value("authorised")
				id   interface{} = r.Context().Value("user_id")
			)
			response.AuthStatus, _ = auth.(bool)
			user_id, ok := id.(int64)
			if ok {
				response.UserId = user_id
			}
			h.APIResponse(w, http.StatusOK, response, "templates/post.html")
		}
	}
}

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == false {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getCreatePost(w, r, "")
	case http.MethodPost:
		h.postCreatePost(w, r)
	default:
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad request"}, "templates/errors.html")
	}
}

func (h *Handler) getCreatePost(w http.ResponseWriter, r *http.Request, errMessage string) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	response := entity.Response{}
	responseChan := make(chan entity.Response)
	go h.forumUcase.FetchCategories(ctx, responseChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case response = <-responseChan:
		err := response.Err
		switch err {
		case entity.ErrInternalServer:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "template/errors.html")
		case nil:
			var (
				auth interface{} = r.Context().Value("authorised")
				id   interface{} = r.Context().Value("user_id")
			)
			response.AuthStatus, _ = auth.(bool)
			user_id, ok := id.(int64)
			if ok {
				response.UserId = user_id
			}
			if errMessage != "" {
				response.ErrorMessage = errMessage
			}
			h.APIResponse(w, http.StatusOK, response, "templates/create_post.html")
		}
	}
}

func (h *Handler) postCreatePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	post, err := entity.GetPost(r)
	if err != nil {
		h.getCreatePost(w, r, err.Error())
		return
	}
	var id interface{} = r.Context().Value("user_id")
	post.User.Id = id.(int64)
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	resChan := make(chan entity.Result)
	var res entity.Result
	go h.forumUcase.StorePost(ctx, post, resChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case res = <-resChan:
		switch res.Err {
		case entity.ErrBadRequest:
			h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"}, "templates/errors.html")
		case nil:
			http.Redirect(w, r, fmt.Sprintf("/posts/%d", res.Id), 303)
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		default:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")

		}
	}
}

func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == false {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	if r.Method != http.MethodPost {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "templates/errors.html")
		return
	}
	r.ParseForm()
	commentRes := entity.GetComment(r)
	if commentRes.Err != nil {
		h.errLog.Println(commentRes.Err)
		if commentRes.Err == entity.ErrEmptyComment {
			h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Empty comment"}, "templates/errors.html")
		} else {
			h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"}, "templates/errors.html")
		}
		return
	}

	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	resChan := make(chan entity.Result)
	var res entity.Result
	go h.forumUcase.StoreComment(ctx, commentRes.Comment, resChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case res = <-resChan:
		switch res.Err {
		case entity.ErrBadRequest:
			h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"}, "templates/errors.html")
		case nil:
			http.Redirect(w, r, fmt.Sprintf("/posts/%d", commentRes.Comment.Post.Id), 303)
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		default:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
		}
	}
}

func (h *Handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "templates/errors.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	response := entity.Response{}
	responseChan := make(chan entity.Response)
	go h.forumUcase.FetchUsers(ctx, responseChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case response = <-responseChan:
		err := response.Err
		switch err {
		case entity.ErrInternalServer:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "template/errors.html")
		case nil:
			var (
				auth interface{} = r.Context().Value("authorised")
				id   interface{} = r.Context().Value("user_id")
			)
			response.AuthStatus, _ = auth.(bool)
			user_id, ok := id.(int64)
			if ok {
				response.UserId = user_id
			}
			h.APIResponse(w, http.StatusOK, response, "templates/users.html")
		}
	}
}

func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "templates/errors.html")
		return
	}
	user_id, err := getID(r.URL.String(), "users")
	if err != nil {
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: err.Error()}, "templates/errors.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	response := entity.Response{}
	responseChan := make(chan entity.Response)
	go h.forumUcase.FetchUser(ctx, user_id, responseChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case response = <-responseChan:
		err := response.Err
		switch err {
		case entity.ErrInternalServer:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		case entity.ErrNotFound:
			h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"}, "templates/errors.html")
		case nil:
			var (
				auth interface{} = r.Context().Value("authorised")
				id   interface{} = r.Context().Value("user_id")
			)
			response.AuthStatus, _ = auth.(bool)
			user_id, ok := id.(int64)
			if ok {
				response.UserId = user_id
			}
			h.APIResponse(w, http.StatusOK, response, "templates/user.html")
		}
	}
}

func (h *Handler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "templates/errors.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	response := entity.Response{}
	responseChan := make(chan entity.Response)
	go h.forumUcase.FetchCategories(ctx, responseChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case response = <-responseChan:
		err := response.Err
		switch err {
		case entity.ErrInternalServer:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "template/errors.html")
		case nil:
			var (
				auth interface{} = r.Context().Value("authorised")
				id   interface{} = r.Context().Value("user_id")
			)
			response.AuthStatus, _ = auth.(bool)
			user_id, ok := id.(int64)
			if ok {
				response.UserId = user_id
			}
			h.APIResponse(w, http.StatusOK, response, "templates/categories.html")
		}
	}
}

func (h *Handler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{ErrorMessage: "Invalid method"}, "templates/errors.html")
		return
	}
	category_id, err := getID(r.URL.String(), "categories")
	if err != nil {
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: err.Error()}, "templates/errors.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	response := entity.Response{}
	responseChan := make(chan entity.Response)
	go h.forumUcase.FetchCategory(ctx, category_id, responseChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case response = <-responseChan:
		err := response.Err
		switch err {
		case entity.ErrInternalServer:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"}, "templates/errors.html")
		case entity.ErrRequestTimeout:
			h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		case entity.ErrNotFound:
			h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"}, "templates/errors.html")
		case nil:
			var (
				auth interface{} = r.Context().Value("authorised")
				id   interface{} = r.Context().Value("user_id")
			)
			response.AuthStatus, _ = auth.(bool)
			user_id, ok := id.(int64)
			if ok {
				response.UserId = user_id
			}
			h.APIResponse(w, http.StatusOK, response, "templates/category.html")
		}
	}
}

func getID(urlPath, endpoint string) (int, error) {
	path := strings.Split(urlPath, "/")
	if path[len(path)-2] != endpoint {
		return 0, errors.New("Not Found")
	}
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		return 0, errors.New("Not Found")
	}
	return id, nil
}

func (h *Handler) PostReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == false {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	r.ParseForm()
	postReaction, err := entity.GetPostReaction(r)
	if err != nil {
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{}, "templates/errors.go")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	errChan := make(chan error)
	go h.forumUcase.PostReaction(ctx, postReaction, errChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case err = <-errChan:
		switch err {
		case nil:
			http.Redirect(w, r, fmt.Sprintf("/posts/%d", postReaction.Post.Id), 303)

		default:
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: err.Error()}, "templates/errors.html")
		}
	}
}

func (h *Handler) CommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("authorised") == false {
		h.APIResponse(w, http.StatusForbidden, entity.Response{ErrorMessage: "Forbidden"}, "templates/errors.html")
		return
	}
	r.ParseForm()
	commentReaction, err := entity.GetCommentReaction(r)
	if err != nil {
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{}, "templates/errors.html")
		return
	}
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	errChan := make(chan error)
	go h.forumUcase.CommentReaction(ctx, commentReaction, errChan)
	select {
	case <-ctx.Done():
		err := ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"}, "templates/errors.html")
		return
	case err = <-errChan:
		switch err {
		case nil:
			http.Redirect(w, r, fmt.Sprintf("/posts/%d#%d", commentReaction.Post.Id, commentReaction.Comment.Id), 303)
		default:
			h.errLog.Println(err)
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: err.Error()}, "templates/errors.html")
		}
	}
}

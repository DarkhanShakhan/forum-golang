package app

import (
	"encoding/json"
	"fmt"
	"forum_app/internal/entity"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) PostDetailsHandler(w http.ResponseWriter, r *http.Request) {
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
	postChan := make(chan entity.PostResult)
	var postResult entity.PostResult
	go h.pcase.FetchById(ctx, id, postChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case postResult = <-postChan:
		err = postResult.Err
		if err != nil {
			if err == entity.ErrPostNotFound {
				h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusOK, entity.Response{Body: postResult.Post})
}
func (h *Handler) PostsAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodGet {
		h.errorLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
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
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case postsRes = <-postsChan:
		if err = postsRes.Err; err != nil {
			h.errorLog.Println(err)
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusOK, entity.Response{Body: postsRes.Posts})
}
func (h *Handler) CategoryPostsHandler(w http.ResponseWriter, r *http.Request) {
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
	catChan := make(chan entity.CatResult)
	var catResult entity.CatResult
	go h.pcase.FetchCategoryPosts(ctx, id, catChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case catResult = <-catChan:
		err = catResult.Err
		if err != nil {
			h.errorLog.Println(err)
			if err == entity.ErrCategoryNotFound {
				h.APIResponse(w, http.StatusNotFound, entity.Response{ErrorMessage: "Not Found"})
				return
			}
			//FIXME:check for existing category
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusOK, entity.Response{Body: catResult.Cat})
}
func (h *Handler) StorePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodPost {
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	var post entity.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil || !validatePostData(post) {
		h.errorLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	resChan := make(chan entity.Result)
	var res entity.Result
	go h.pcase.Store(ctx, post, resChan)
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
			if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusCreated, entity.Response{Body: entity.Post{Id: int(res.Id)}})
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
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodPost {
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	var post_reaction entity.PostReaction
	err := json.NewDecoder(r.Body).Decode(&post_reaction)
	if err != nil || !validatePostReactionData(post_reaction) {
		h.errorLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{})
		return
	}
	errChan := make(chan error)
	go h.pcase.StorePostReaction(ctx, post_reaction, errChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case err = <-errChan:
		if err != nil {
			h.errorLog.Println(err)
			if isConstraintError(err) {
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusNoContent, entity.Response{})
}

func validatePostReactionData(reaction entity.PostReaction) bool {
	return reaction.Post.Id != 0 && reaction.Reaction.User.Id != 0
}

func isConstraintError(err error) bool {
	return strings.Contains(err.Error(), "FOREIGN KEY constraint failed") || strings.Contains(err.Error(), "UNIQUE constraint failed: post_reactions.post_id, post_reactions.user_id")
}

func (h *Handler) UpdatePostReactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodPut {
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	var post_reaction entity.PostReaction
	err := json.NewDecoder(r.Body).Decode(&post_reaction)
	if err != nil || validatePostReactionData(post_reaction) {
		h.errorLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	errChan := make(chan error)
	go h.pcase.UpdatePostReaction(ctx, post_reaction, errChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case err = <-errChan:
		if err != nil {
			h.errorLog.Println(err)
			if isConstraintError(err) {
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusNoContent, entity.Response{})
}

func (h *Handler) DeletePostReactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodDelete {
		h.errorLog.Printf("invalid method: %s\n", r.Method)
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	var post_reaction entity.PostReaction
	err := json.NewDecoder(r.Body).Decode(&post_reaction)
	if err != nil || !validatePostReactionData(post_reaction) {
		h.errorLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{})
		return
	}
	errChan := make(chan error)
	go h.pcase.DeletePostReaction(ctx, post_reaction, errChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errorLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return

	case err = <-errChan:
		if err != nil {
			h.errorLog.Println(err)
			if isConstraintError(err) {
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusNoContent, entity.Response{})
}

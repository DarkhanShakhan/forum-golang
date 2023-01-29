package app

import (
	"encoding/json"
	"fmt"
	"forum_app/internal/entity"
	"net/http"
	"strconv"
)

func (h *Handler) StoreCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodPost {
		h.errLog.Printf("invalid method: %s\n", r.Method)
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	var comment entity.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil || !validateCommentData(comment) {
		h.errLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	resChan := make(chan entity.Result)
	var result entity.Result
	go h.ccase.Store(ctx, comment, resChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case result = <-resChan:
		if result.Err != nil {
			h.errLog.Println(err)
			if isConstraintError(err) {
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusCreated, entity.Response{Body: entity.Comment{Id: int(result.Id)}})
}

func validateCommentData(comment entity.Comment) bool {
	// FIXME:check for empty comment
	return comment.Content != "" && comment.Post.Id != 0 && comment.User.Id != 0
}

func validateCommentReactionData(reaction entity.CommentReaction) bool {
	return reaction.Comment.Id != 0 && reaction.Reaction.User.Id != 0
}

func (h *Handler) StoreCommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodPost {
		h.errLog.Printf("invalid method: %s\n", r.Method)
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	var comment_reaction entity.CommentReaction
	err := json.NewDecoder(r.Body).Decode(&comment_reaction)
	// FIXME:validate data
	if err != nil || !validateCommentReactionData(comment_reaction) {
		h.errLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{})
		return
	}
	errChan := make(chan error)
	go h.ccase.StoreCommentReaction(ctx, comment_reaction, errChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case err = <-errChan:
		if err != nil {
			h.errLog.Println(err)
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

func (h *Handler) UpdateCommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodPut {
		h.errLog.Printf("invalid method: %s\n", r.Method)
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	var comment_reaction entity.CommentReaction
	err := json.NewDecoder(r.Body).Decode(&comment_reaction)
	if err != nil || !validateCommentReactionData(comment_reaction) {
		fmt.Println(err)
		fmt.Println(comment_reaction)
		h.errLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	errChan := make(chan error)
	go h.ccase.UpdateCommentReaction(ctx, comment_reaction, errChan)
	select {
	case err = <-errChan:
		if err != nil {
			h.errLog.Println(err)
			if isConstraintError(err) || isNoRowAffectedError(err) {
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	case <-ctx.Done():
		err = ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	}
	h.APIResponse(w, http.StatusNoContent, entity.Response{})
}

func (h *Handler) DeleteCommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodDelete {
		h.errLog.Printf("invalid method: %s\n", r.Method)
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	var comment_reaction entity.CommentReaction
	err := json.NewDecoder(r.Body).Decode(&comment_reaction)
	if err != nil || !validateCommentReactionData(comment_reaction) {
		h.errLog.Println("bad request")
		h.APIResponse(w, http.StatusBadRequest, entity.Response{})
		return
	}
	errChan := make(chan error)
	go h.ccase.DeleteCommentReaction(ctx, comment_reaction, errChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case err = <-errChan:
		if err != nil {
			h.errLog.Println(err)
			if isConstraintError(err) || isNoRowAffectedError(err) {
				h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
				return
			}
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusNoContent, entity.Response{})
}

func (h *Handler) CommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getTimeout(r.Context())
	defer cancel()
	if r.Method != http.MethodGet {
		h.errLog.Println(fmt.Sprintf("method not allowed: %s", r.Method))
		h.APIResponse(w, http.StatusMethodNotAllowed, entity.Response{})
		return
	}
	if err := r.ParseForm(); err != nil {
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusBadRequest, entity.Response{ErrorMessage: "Bad Request"})
		return
	}
	reactionsChan := make(chan entity.ReactionsResult)
	var reactionsRes entity.ReactionsResult
	go h.ccase.FetchReactions(ctx, id, reactionsChan)
	select {
	case <-ctx.Done():
		err = ctx.Err()
		h.errLog.Println(err)
		h.APIResponse(w, http.StatusRequestTimeout, entity.Response{ErrorMessage: "Request Timeout"})
		return
	case reactionsRes = <-reactionsChan:
		if err = reactionsRes.Err; err != nil {
			h.errLog.Println(err)
			h.APIResponse(w, http.StatusInternalServerError, entity.Response{ErrorMessage: "Internal Server Error"})
			return
		}
	}
	h.APIResponse(w, http.StatusOK, entity.Response{Body: reactionsRes.Reactions})
}

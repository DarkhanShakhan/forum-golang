package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"forum_gateway/internal/entity"
	"log"
	"net/http"
)

type ForumUsecase struct {
	errLog *log.Logger
}

func NewForumUsecase(errLog *log.Logger) *ForumUsecase {
	return &ForumUsecase{errLog: errLog}
}

func (f *ForumUsecase) FetchPosts(ctx context.Context, responseChan chan entity.Response) {
	response, err := getAPIResponse(ctx, http.MethodGet, "http://localhost:8080/posts", []byte{})
	if err != nil {
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
		return
	}
	switch response.StatusCode {
	case 408:
		responseChan <- entity.Response{Err: entity.ErrRequestTimeout}
	case 200:
		result, err := getResponse(response.Body)
		if err != nil {
			responseChan <- entity.Response{Err: entity.ErrInternalServer}
			return
		}
		responseChan <- result
	default:
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
	}
}

func (f *ForumUsecase) FetchPost(ctx context.Context, id int, responseChan chan entity.Response) {
	response, err := getAPIResponse(ctx, http.MethodGet, fmt.Sprintf("http://localhost:8080/post?id=%d", id), []byte{})
	if err != nil {
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
		return
	}
	switch response.StatusCode {
	case 408:
		responseChan <- entity.Response{Err: entity.ErrRequestTimeout}
	case 200:
		result, err := getResponse(response.Body)
		if err != nil {
			responseChan <- entity.Response{Err: entity.ErrInternalServer}
			return
		}
		responseChan <- result
	case 404:
		responseChan <- entity.Response{Err: entity.ErrNotFound}
	default:
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
	}
}

func (f *ForumUsecase) StorePost(ctx context.Context, post entity.Post, resChan chan entity.Result) {
	body, err := json.Marshal(post)
	if err != nil {
		resChan <- entity.Result{Err: entity.ErrInternalServer}
		return
	}
	response, err := getAPIResponse(ctx, http.MethodPost, "http://localhost:8080/post/save", body)
	if err != nil {
		resChan <- entity.Result{Err: entity.ErrInternalServer}
		return
	}
	switch response.StatusCode {
	case 408:
		resChan <- entity.Result{Err: entity.ErrRequestTimeout}
	case 400:
		resChan <- entity.Result{Err: entity.ErrBadRequest}
	case 201:
		resPost := getPost(response.Body)
		if resPost.Err != nil {
			f.errLog.Println(resPost.Err)
			resChan <- entity.Result{Err: entity.ErrInternalServer}
			return
		}
		resChan <- entity.Result{Id: resPost.Post.Id}
	default:
		resChan <- entity.Result{Err: entity.ErrInternalServer}
	}
}

func (f *ForumUsecase) StoreComment(ctx context.Context, comment entity.Comment, resChan chan entity.Result) {
	body, err := json.Marshal(comment)
	if err != nil {
		resChan <- entity.Result{Err: entity.ErrInternalServer}
		return
	}
	response, err := getAPIResponse(ctx, http.MethodPost, "http://localhost:8080/comment/save", body)
	if err != nil {
		resChan <- entity.Result{Err: entity.ErrInternalServer}
		return
	}
	switch response.StatusCode {
	case 408:
		resChan <- entity.Result{Err: entity.ErrRequestTimeout}
	case 400:
		resChan <- entity.Result{Err: entity.ErrBadRequest}
	case 201:
		resComment := getComment(response.Body)
		if resComment.Err != nil {
			f.errLog.Println(resComment.Err)
			resChan <- entity.Result{Err: entity.ErrInternalServer}
			return
		}
		resChan <- entity.Result{Id: resComment.Comment.Id}
	default:
		resChan <- entity.Result{Err: entity.ErrInternalServer}
	}
}

func (f *ForumUsecase) FetchUsers(ctx context.Context, responseChan chan entity.Response) {
	response, err := getAPIResponse(ctx, http.MethodGet, "http://localhost:8080/users", []byte{})
	if err != nil {
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
		return
	}
	switch response.StatusCode {
	case 408:
		responseChan <- entity.Response{Err: entity.ErrRequestTimeout}
	case 200:
		result, err := getResponse(response.Body)
		if err != nil {
			responseChan <- entity.Response{Err: entity.ErrInternalServer}
			return
		}
		responseChan <- result
	default:
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
	}
}

func (f *ForumUsecase) FetchUser(ctx context.Context, id int, responseChan chan entity.Response) {
	response, err := getAPIResponse(ctx, http.MethodGet, fmt.Sprintf("http://localhost:8080/user?id=%d", id), []byte{})
	if err != nil {
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
		return
	}
	switch response.StatusCode {
	case 408:
		responseChan <- entity.Response{Err: entity.ErrRequestTimeout}
	case 200:
		result, err := getResponse(response.Body)
		if err != nil {
			responseChan <- entity.Response{Err: entity.ErrInternalServer}
			return
		}
		responseChan <- result
	case 404:
		responseChan <- entity.Response{Err: entity.ErrNotFound}
	default:
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
	}
}

func (f *ForumUsecase) FetchCategory(ctx context.Context, id int, responseChan chan entity.Response) {
	response, err := getAPIResponse(ctx, http.MethodGet, fmt.Sprintf("http://localhost:8080/category?id=%d", id), []byte{})
	if err != nil {
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
		return
	}
	switch response.StatusCode {
	case 408:
		responseChan <- entity.Response{Err: entity.ErrRequestTimeout}
	case 200:
		result, err := getResponse(response.Body)
		if err != nil {
			responseChan <- entity.Response{Err: entity.ErrInternalServer}
			return
		}
		responseChan <- result
	case 404:
		responseChan <- entity.Response{Err: entity.ErrNotFound}
	default:
		responseChan <- entity.Response{Err: entity.ErrInternalServer}
	}
}

func (f *ForumUsecase) PostReaction(ctx context.Context, reaction entity.PostReaction, errorChan chan error) {
	response, _ := getAPIResponse(ctx, http.MethodGet, fmt.Sprintf("http://localhost:8080/post_reactions?id=%d", reaction.Post.Id), nil)
	res := getReactions(response.Body)
	for _, i := range res.Reactions {
		if i.User.Id == reaction.Reaction.User.Id {
			fmt.Println("here")
		}
	}
	errorChan <- nil
}

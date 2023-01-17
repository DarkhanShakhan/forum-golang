package usecase

import (
	"context"
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

func (f *ForumUsecase) FetchPost(id int, response chan entity.Response) {
}

func (f *ForumUsecase) FetchUsers(response chan entity.Response) {
}

func (f *ForumUsecase) FetchUser(id int, response chan entity.Response) {
}

func (f *ForumUsecase) FetchCategory(response chan entity.Response) {
}

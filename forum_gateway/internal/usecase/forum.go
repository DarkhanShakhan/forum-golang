package usecase

import (
	"forum_gateway/internal/entity"
	"log"
)

type ForumUsecase struct {
	errLog *log.Logger
}

func (f *ForumUsecase) FetchPosts(response chan entity.Response) {
}

func (f *ForumUsecase) FetchPost(id int, response chan entity.Response) {
}

func (f *ForumUsecase) FetchUsers(response chan entity.Response) {
}

func (f *ForumUsecase) FetchUser(id int, response chan entity.Response) {
}

func (f *ForumUsecase) FetchCategory(response chan entity.Response) {
}

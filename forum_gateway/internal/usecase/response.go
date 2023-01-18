package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"forum_gateway/internal/entity"
	"io"
	"net/http"
)

func getAPIResponse(ctx context.Context, method string, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	return client.Do(req)
}

func getResponse(response io.ReadCloser) (entity.Response, error) {
	result := entity.Response{}
	err := json.NewDecoder(response).Decode(&result)
	return result, err
}

func getSession(response io.ReadCloser) (entity.Session, error) { // FIXME: change to SessionRes
	temp, err := getResponse(response)
	if err != nil {
		return entity.Session{}, err
	}
	jsonSession, err := json.Marshal(temp.Body)
	if err != nil {
		return entity.Session{}, err
	}
	session := entity.Session{}
	err = json.Unmarshal(jsonSession, &session)
	return session, err
}

func getPost(response io.ReadCloser) entity.PostResult {
	temp, err := getResponse(response)
	if err != nil {
		return entity.PostResult{Err: err}
	}
	jsonPost, err := json.Marshal(temp.Body)
	if err != nil {
		return entity.PostResult{Err: err}
	}
	post := entity.Post{}
	err = json.Unmarshal(jsonPost, &post)
	return entity.PostResult{Err: err, Post: post}
}

func getAuthStatus(response io.ReadCloser) entity.AuthStatusResult {
	temp, err := getResponse(response)
	if err != nil {
		return entity.AuthStatusResult{Err: err}
	}
	jsonAuthStatus, err := json.Marshal(temp.Body)
	if err != nil {
		return entity.AuthStatusResult{Err: err}
	}
	authStatus := entity.AuthStatusResult{}
	err = json.Unmarshal(jsonAuthStatus, &authStatus)
	if err != nil {
		return entity.AuthStatusResult{Err: err}
	}
	return authStatus
}

func getComment(response io.ReadCloser) entity.CommentResult {
	temp, err := getResponse(response)
	if err != nil {
		return entity.CommentResult{Err: err}
	}
	jsonComment, err := json.Marshal(temp.Body)
	if err != nil {
		return entity.CommentResult{Err: err}
	}
	comment := entity.Comment{}
	err = json.Unmarshal(jsonComment, &comment)
	return entity.CommentResult{Err: err, Comment: comment}
}

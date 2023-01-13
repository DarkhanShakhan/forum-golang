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

func getSession(response io.ReadCloser) (entity.Session, error) {
	temp := entity.Response{}
	err := json.NewDecoder(response).Decode(&temp)
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

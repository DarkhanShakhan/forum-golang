package app

import (
	"bytes"
	"context"
	"forum_gateway/internal/entity"
	"html/template"
	"net/http"
)

func (h *Handler) APIResponse(w http.ResponseWriter, code int, response entity.Response, filename string) {
	templ, err := template.ParseFiles(filename)
	if err != nil {
		h.errLog.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	templ.Execute(w, response)
}

func getAPIResponse(ctx context.Context, method string, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	return client.Do(req)
}

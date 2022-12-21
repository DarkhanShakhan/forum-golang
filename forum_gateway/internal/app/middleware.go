package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUrl := "http://localhost:8081/authenticate"
		cookier, err := r.Cookie("token")
		if err != nil {
			ctx := context.WithValue(r.Context(), "authorised", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		jsonStr := []byte(fmt.Sprintf(`{"token":"%s"}`, cookier.Value))
		req, err := http.NewRequest("GET", requestUrl, bytes.NewBuffer(jsonStr))
		if err != nil {
			ctx := context.WithValue(r.Context(), "authorised", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ctx := context.WithValue(r.Context(), "authorised", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		var authStatus AuthStatusResult
		if err = json.NewDecoder(resp.Body).Decode(&authStatus); err != nil {
			ctx := context.WithValue(r.Context(), "authorised", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		if authStatus.Status == NonAuthorised {
			ctx := context.WithValue(r.Context(), "authorised", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		cookie := http.Cookie{
			Name:  "token",
			Value: authStatus.Token,
			Path:  "/",
		}
		http.SetCookie(w, &cookie)

		// FIXME: validate data
		ctx := context.WithValue(context.WithValue(r.Context(), "authorised", true), "user_id", authStatus.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type AuthStatus int

const (
	NonAuthorised AuthStatus = iota
	Authorised
)

type AuthStatusResult struct {
	Status AuthStatus `json:"status,omitempty"`
	UserId int64      `json:"user_id,omitempty"`
	Token  string     `json:"token,omitempty"`
	Err    error      `json:"error,omitempty"`
}

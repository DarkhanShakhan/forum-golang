package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"forum_gateway/internal/entity"
	"net/http"
)

func (h *Handler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookier, err := r.Cookie("token")
		if err != nil || cookier.Value == "" {
			ctx := context.WithValue(r.Context(), "authorised", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx, cancel := getTimeout(r.Context())
		defer cancel()

		authResChan := make(chan entity.AuthStatusResult)
		go h.auUcase.Authenticate(ctx, cookier.Value, authResChan)
		select {
		case authRes := <-authResChan:
			if authRes.Status == entity.Authorised {
				cookie := http.Cookie{
					Name:  "token",
					Value: authRes.Session.Token,
					Path:  "/",
				}
				http.SetCookie(w, &cookie)
				ctx := context.WithValue(context.WithValue(r.Context(), "authorised", true), "user_id", authRes.Session.UserId)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				h.errLog.Println(authRes.Err)
				next.ServeHTTP(w, r)
			}
		case <-ctx.Done():
			err = ctx.Err()
			h.errLog.Println(err)
			next.ServeHTTP(w, r)
			return
		}
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUrl := "http://localhost:8081/authenticate"
		cookier, err := r.Cookie("token")
		fmt.Println(cookier)
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
		//	FIXME: expiry time
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

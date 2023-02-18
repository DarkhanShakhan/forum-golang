package app

import (
	"context"
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
				if authRes.Err != nil {
					h.errLog.Println(authRes.Err)
				}
				ctx := context.WithValue(r.Context(), "authorised", false)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		case <-ctx.Done():
			err = ctx.Err()
			h.errLog.Println(err)
			ctx := context.WithValue(r.Context(), "authorised", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	})
}

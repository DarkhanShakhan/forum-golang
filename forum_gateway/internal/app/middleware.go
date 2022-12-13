package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUrl := "http://localhost:8081/authenticate"
		cookier, err := r.Cookie("token")
		if err != nil {
			fmt.Println(err)
			return
		}
		jsonStr := []byte(fmt.Sprintf(`{"token":"%s"}`, cookier.Value))
		req, err := http.NewRequest("GET", requestUrl, bytes.NewBuffer(jsonStr))

		client := &http.Client{}
		resp, err := client.Do(req)
		fmt.Println(err)
		var authStatus AuthStatusResult
		json.NewDecoder(resp.Body).Decode(&authStatus) // error
		cookie := http.Cookie{
			Name:  "token",
			Value: authStatus.Token,
		}
		http.SetCookie(w, &cookie)
		fmt.Println(authStatus)
	})
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	requestUrl := "http://localhost:8081/auth"
	// 	cookie, err := r.Cookie("token")
	// 	rr := r.Cookies()
	// 	fmt.Println(rr)
	// 	fmt.Println(err)
	// 	if err != nil {
	// 		return
	// 	}
	// 	// FIXME: deal err
	// 	jsonStr := []byte(fmt.Sprintf(`{"token":"%s"}`, cookie.Value))
	// 	req, err := http.NewRequest("GET", requestUrl, bytes.NewBuffer(jsonStr))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	req.Header.Set("Content-Type", "application/json")

	// 	client := &http.Client{}
	// 	resp, _ := client.Do(req)
	// 	// defer resp.Body.Close()
	// 	var authStatus AuthStatusResult
	// 	json.NewDecoder(resp.Body).Decode(&authStatus) // error
	// 	fmt.Println(authStatus)
	// }
}

type AuthStatus int

const (
	NonAuthorised AuthStatus = iota
	Authorised
)

type AuthStatusResult struct {
	Status AuthStatus `json:"status,omitempty"`
	Token  string     `json:"token,omitempty"`
	Err    error      `json:"error,omitempty"`
}

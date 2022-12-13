package app

import (
	"fmt"
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/sign_in", f)
	handler := Authenticate(mux)
	srv := &http.Server{
		Addr:    "localhost:8082",
		Handler: handler,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}

func f(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
}

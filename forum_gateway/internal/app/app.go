package app

import (
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	mux.Handle("/signin_google", Authenticate(http.HandlerFunc(SignInGoogleHandler)))
	mux.HandleFunc("/google_callback", GoogleCallbackHandler)
	// handler := Authenticate(mux)
	srv := &http.Server{
		Addr:    "localhost:8082",
		Handler: mux,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}

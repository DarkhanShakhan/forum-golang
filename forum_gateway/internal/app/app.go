package app

import (
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	mux.Handle("/signin_google", Authenticate(http.HandlerFunc(SignInGoogleHandler)))
	mux.HandleFunc("/google_callback", GoogleCallbackHandler)
	mux.Handle("/posts", Authenticate(http.HandlerFunc(PostsHandler)))
	mux.Handle("/post/", Authenticate(http.HandlerFunc(PostHandler)))
	mux.Handle("/users", Authenticate(http.HandlerFunc(UsersHandler)))
	mux.Handle("/user/", Authenticate(http.HandlerFunc(UserHandler)))
	mux.Handle("/category/", Authenticate(http.HandlerFunc(CategoryHandler)))
	mux.Handle("/post/create", Authenticate(http.HandlerFunc(PostCreateHandler)))
	mux.Handle("/sign_in", Authenticate(http.HandlerFunc(SignInHandler)))

	// handler := Authenticate(mux)
	srv := &http.Server{
		Addr:    "localhost:8082",
		Handler: mux,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}

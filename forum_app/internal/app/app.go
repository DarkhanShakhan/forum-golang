package app

import (
	"log"
	"net/http"
	"os"
)

func Run() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	h := NewHandler(errorLog)
	mux := http.NewServeMux()
	mux.HandleFunc("/users", h.UsersAllHandler)
	mux.HandleFunc("/user/id", h.UserByIdHandler)
	mux.HandleFunc("/user/email", h.UserByEmailHandler)
	mux.HandleFunc("/post", h.PostFullHandler)
	mux.HandleFunc("/posts", h.PostsAllHandler)
	mux.HandleFunc("/post/save", h.StorePostHandler)

	srv := &http.Server{
		Addr:     "localhost:8080",
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Println("Listening on localhost:8080")
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

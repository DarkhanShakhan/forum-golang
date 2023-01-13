package app

import (
	"io"
	"log"
	"net/http"
	"os"
)

// FIXME: add info logging when database open
// TODO: middleware for inserting context deadline

func Run() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		errLog.Println("Log file doesn't open")
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stderr, f)
	errLog.SetOutput(wrt)
	infoLog.SetOutput(wrt)
	h := NewHandler(errLog, infoLog)
	mux := http.NewServeMux()
	// get
	mux.HandleFunc("/users", h.UsersAllHandler)
	mux.HandleFunc("/user", h.UserDetailsHandler)
	mux.HandleFunc("/user/email", h.UserByEmailHandler)
	mux.HandleFunc("/post", h.PostDetailsHandler)
	mux.HandleFunc("/posts", h.PostsAllHandler)
	mux.HandleFunc("/category", h.CategoryPostsHandler)

	// post
	mux.HandleFunc("/user/save", h.StoreUserHandler)
	mux.HandleFunc("/post/save", h.StorePostHandler)
	mux.HandleFunc("/post_reaction/save", h.StorePostReactionHandler)
	mux.HandleFunc("/comment/save", h.StoreCommentHandler)
	mux.HandleFunc("/comment_reaction/save", h.StoreCommentReactionHandler)

	// put
	mux.HandleFunc("/post_reaction/update", h.UpdatePostReactionHandler)
	mux.HandleFunc("/comment_reaction/update", h.UpdateCommentReactionHandler)

	// delete
	mux.HandleFunc("/post_reaction/delete", h.DeletePostReactionHandler)
	srv := &http.Server{
		Addr:     "localhost:8080",
		ErrorLog: errLog,
		Handler:  mux,
	}

	infoLog.Println("Listening on localhost:8080")
	err = srv.ListenAndServe()
	errLog.Fatal(err)
	// TODO: graceful shutdown
}

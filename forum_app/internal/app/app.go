package app

import (
	"io"
	"log"
	"net/http"
	"os"
)

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
	mux.HandleFunc("/post_reactions", h.PostReactionsHandler)
	mux.HandleFunc("/comment_reactions", h.CommentReactionHandler)
	mux.HandleFunc("/category", h.CategoryPostsHandler)
	mux.HandleFunc("/categories", h.CategoriesHandler)

	// post
	mux.HandleFunc("/user/save", h.StoreUserHandler)
	mux.HandleFunc("/post/save", h.StorePostHandler)
	mux.HandleFunc("/post_reactions/save", h.StorePostReactionHandler)
	mux.HandleFunc("/comments/save", h.StoreCommentHandler)
	mux.HandleFunc("/comment_reactions/save", h.StoreCommentReactionHandler)

	// put
	mux.HandleFunc("/post_reactions/update", h.UpdatePostReactionHandler)
	mux.HandleFunc("/comment_reactions/update", h.UpdateCommentReactionHandler)

	// delete
	mux.HandleFunc("/comment_reactions/delete", h.DeleteCommentReactionHandler)
	mux.HandleFunc("/post_reactions/delete", h.DeletePostReactionHandler)
	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: errLog,
		Handler:  mux,
	}

	infoLog.Println("Listening on localhost:8080")
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}

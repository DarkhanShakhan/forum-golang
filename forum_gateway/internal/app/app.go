package app

import (
	"forum_gateway/internal/usecase"
	"io"
	"log"
	"net/http"
	"os"
)

func Run() {
	mux := http.NewServeMux()
	infoLog, errLog, file := getLogs("log.txt")
	defer file.Close()
	auUcase := usecase.NewAuthUsecase(errLog, infoLog)
	forumUcase := usecase.NewForumUsecase(errLog)
	h := NewHandler(errLog, infoLog, auUcase, forumUcase)
	// auth
	mux.Handle("/sign-up", h.MultipleMiddleware(h.SignUpHandler))
	mux.Handle("/sign-in", h.MultipleMiddleware(h.SignInHandler))
	mux.Handle("/sign-out", h.MultipleMiddleware(h.SignOutHandler))

	mux.Handle("/sign-in/google", h.MultipleMiddleware(h.SignInOAuthHandler))
	mux.HandleFunc("/callback/google", h.CallBackHandler)

	mux.Handle("/sign-in/github", h.MultipleMiddleware(h.SignInOAuthHandler))
	mux.HandleFunc("/callback/github", h.CallBackHandler)

	// forum
	mux.Handle("/", h.MultipleMiddleware(h.PostsHandler))
	mux.Handle("/posts", h.MultipleMiddleware(h.PostsHandler))
	mux.Handle("/posts/", h.MultipleMiddleware(h.PostHandler))
	mux.Handle("/posts/new", h.MultipleMiddleware(h.CreatePostHandler))
	mux.Handle("/comments/new", h.MultipleMiddleware(h.CreateCommentHandler))
	mux.Handle("/users", h.MultipleMiddleware(h.UsersHandler))
	mux.Handle("/users/", h.MultipleMiddleware(h.UserHandler))
	mux.Handle("/categories", h.MultipleMiddleware(h.CategoriesHandler))
	mux.Handle("/categories/", h.MultipleMiddleware(h.CategoryHandler))
	mux.Handle("/post-reactions/new", h.MultipleMiddleware(h.PostReactionHandler))
	mux.Handle("/comment-reactions/new", h.MultipleMiddleware(h.CommentReactionHandler))

	mux.Handle("/templates/css/", http.StripPrefix("/templates/css/", http.FileServer(http.Dir("templates/css"))))
	mux.Handle("/templates/img/", http.StripPrefix("/templates/img/", http.FileServer(http.Dir("templates/img"))))

	srv := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}
	infoLog.Println("Listening on localhost:8082")
	err := srv.ListenAndServe()
	log.Fatal(err)
}

func getLogs(filename string) (*log.Logger, *log.Logger, *os.File) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		errorLog.Println("Log file doesn't open")
	}
	wrt := io.MultiWriter(os.Stderr, f)
	errorLog.SetOutput(wrt)
	infoLog.SetOutput(wrt)
	return infoLog, errorLog, f
}

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
	h := NewHandler(errLog, infoLog, auUcase)
	// auth
	mux.Handle("/sign_up", Authenticate(http.HandlerFunc(h.SignUpHandler)))

	mux.Handle("/signin_google", Authenticate(http.HandlerFunc(SignInGoogleHandler)))
	mux.HandleFunc("/google_callback", GoogleCallbackHandler)
	mux.Handle("/posts", Authenticate(http.HandlerFunc(PostsHandler)))
	// mux.Handle("/posts/", Authenticate(http.HandlerFunc(PostHandler)))
	mux.Handle("/users", Authenticate(http.HandlerFunc(UsersHandler)))
	mux.Handle("/users/", Authenticate(http.HandlerFunc(UserHandler)))
	mux.Handle("/categories/", Authenticate(http.HandlerFunc(CategoryHandler)))
	mux.Handle("/posts/new", Authenticate(http.HandlerFunc(PostCreateHandler)))
	mux.Handle("/sign_in", Authenticate(http.HandlerFunc(h.SignInHandler)))

	mux.Handle("/sign_out", Authenticate(http.HandlerFunc(SignOutHandler)))
	mux.Handle("/comments/new", Authenticate(http.HandlerFunc(CommentCreateHandler)))

	// handler := Authenticate(mux)
	srv := &http.Server{
		Addr:    "localhost:8082",
		Handler: mux,
	}

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

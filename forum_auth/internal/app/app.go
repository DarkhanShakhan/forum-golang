package app

import (
	"io"
	"log"
	"net/http"
	"os"
)

func Run() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	wrt := io.MultiWriter(os.Stderr, f)
	errorLog.SetOutput(wrt)
	h := NewHandler(errorLog)
	mux := http.NewServeMux()

	mux.HandleFunc("/sign_in", h.SignInHandler)
	mux.HandleFunc("/sign_out", h.SignOutHandler)
	mux.HandleFunc("/authenticate", h.Authenticate)
	mux.HandleFunc("/sign_up", h.SignUpHandler)
	srv := &http.Server{
		Addr:     "localhost:8081",
		ErrorLog: errorLog,
		Handler:  mux,
	}
	infoLog.Println("Listening on localhost:8081")
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
	// TODO: graceful shutdown
}

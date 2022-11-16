package app

import (
	cr "forum/internal/forum_app/comment/repository"
	"log"
	"net/http"

	cUcse "forum/internal/forum_app/comment/usecase"
	pr "forum/internal/forum_app/post/repository"
	pUcse "forum/internal/forum_app/post/usecase"
	ur "forum/internal/forum_app/user/repository"
	uUcse "forum/internal/forum_app/user/usecase"

	"forum/pkg/sqlite3"
)

func Run() {
	db, _ := sqlite3.New()
	usersRepo := ur.NewUsersRepository(db)
	postsRepo := pr.NewPostsRepository(db)
	pReactionsRepo := pr.NewPostReactionsRepository(db)
	categoriesRepo := pr.NewCategoriesRepository(db)
	commentsRepo := cr.NewCommentsRepository(db)
	cReactionsRepo := cr.NewCommentReactionsRepository(db)
	ucase := uUcse.NewUsersUsecase(usersRepo, postsRepo, pReactionsRepo, commentsRepo, cReactionsRepo)
	pcase := pUcse.NewPostsUsecase(postsRepo, pReactionsRepo, commentsRepo, categoriesRepo, usersRepo)
	ccase := cUcse.NewCommentsUsecase(commentsRepo, cReactionsRepo, postsRepo, usersRepo)

	server := NewServer(ucase, pcase, ccase)
	http.HandleFunc("/users", server.UsersAllHandler)
	http.HandleFunc("/user", server.UserByIdHandler)
	http.HandleFunc("/post", server.PostByIdHandler)
	http.HandleFunc("/posts", server.PostsAllHandler)
	http.HandleFunc("/post/save", server.StorePostHandler)
	log.Print("Listening...")
	http.ListenAndServe(":8080", nil)
}

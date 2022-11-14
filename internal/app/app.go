package app

import (
	cr "forum/internal/comment/repository"
	"forum/internal/gateway"
	"log"
	"net/http"

	cUcse "forum/internal/comment/usecase"
	pr "forum/internal/post/repository"
	pUcse "forum/internal/post/usecase"
	ur "forum/internal/user/repository"
	uUcse "forum/internal/user/usecase"

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

	apiGateway := gateway.NewAPIGateway(ucase, pcase, ccase)
	http.HandleFunc("/main", apiGateway.MainHandler)
	http.HandleFunc("/user/", apiGateway.UserHandler)
	log.Print("Listening...")
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	cr "forum/internal/comment/repository"

	// pr "forum/internal/post/repository"
	// ur "forum/internal/user/repository"

	"forum/pkg/sqlite3"
)

func main() {
	db, _ := sqlite3.New()
	// usersRepo := ur.NewUsersRepository(db)
	// postsRepo := pr.NewPostsRepository(db)
	// pReactionsRepo := pr.NewPostReactionsRepository(db)
	// categoriesRepo := pr.NewCategoriesRepository(db)
	commentsRepo := cr.NewCommentsRepository(db)
	// user1, _ := usersRepo.FetchById(2)
	// post2, _ := postsRepo.FetchById(1)
	// comment1 := entity.Comment{Post: post2, User: user1, Content: "Hello W"}
	// fmt.Println(comment1)
	fmt.Println(commentsRepo.FetchByUserId(1))
}

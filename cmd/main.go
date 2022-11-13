package main

import (
	"fmt"
	"forum/internal/entity"
	pr "forum/internal/post/repository"

	ur "forum/internal/user/repository"
	"forum/pkg/sqlite3"
)

func main() {
	db, _ := sqlite3.New()
	usersRepo := ur.NewUsersRepository(db)
	postsRepo := pr.NewPostsRepository(db)
	pReactionsRepo := pr.NewPostReactionsRepository(db)
	// post1 := entity.Post{User: entity.User{Id: 1}, Title: "post1"}
	// post2 := entity.Post{User: entity.User{Id: 2}, Title: "post2"}
	// postsRepo.Store(post1)
	// postsRepo.Store(post2)
	user1, _ := usersRepo.FetchById(1)
	post2, _ := postsRepo.FetchById(2)
	fmt.Println(user1)
	fmt.Println(post2)
	postReaction1 := entity.PostReaction{Reaction: entity.Reaction{User: user1, Like: true}, Post: post2}
	fmt.Println(pReactionsRepo.StoreReaction(postReaction1))
}

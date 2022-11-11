package main

import (
	"fmt"
	pr "forum/internal/post/repository"
	"forum/pkg/sqlite3"
)

func main() {
	db, _ := sqlite3.New()
	// usersRepo := ur.NewUsersRepository(db)
	postsRepo := pr.NewPostsRepository(db)
	fmt.Println(postsRepo.FetchById(13))

}

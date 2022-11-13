package main

import (

	// ur "forum/internal/user/repository"

	"forum/pkg/sqlite3"
)

func main() {
	sqlite3.New()
	// usersRepo := ur.NewUsersRepository(db)
	// postsRepo := pr.NewPostsRepository(db)
	// pReactionsRepo := pr.NewPostReactionsRepository(db)
	// categoriesRepo := pr.NewCategoriesRepository(db)

}

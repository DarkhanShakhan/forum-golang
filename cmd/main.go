package main

import (
	"fmt"
	ur "forum/internal/user/repository"
	"forum/pkg/sqlite3"
)

func main() {
	db, _ := sqlite3.New()
	usersRepo := ur.NewUsersRepository(db)
	user, _ := usersRepo.FetchById(1)
	fmt.Println(user)
}

package entity

import "errors"

var (
	ErrUserNotFound     = errors.New("user doesn't exist")
	ErrUserExists       = errors.New("user with a given email already exists")
	ErrPostNotFound     = errors.New("post doesn't exist")
	ErrCategoryNotFound = errors.New("category doesn't exist")
)

package entity

import "errors"

var (
	ErrRequestTimeout  = errors.New("Request Timeout")
	ErrNotFound        = errors.New("Not Found")
	ErrInternalServer  = errors.New("Internal Server Error")
	ErrInvalidPassword = errors.New("Invalid password")
	ErrEmailExists     = errors.New("Email already exists")
	ErrBadRequest      = errors.New("Bad Request")
)

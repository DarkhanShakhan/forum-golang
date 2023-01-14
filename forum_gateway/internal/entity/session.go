package entity

import (
	"time"
)

type Session struct {
	Token      string    `json:"token,omitempty"`
	UserId     int64     `json:"user_id,omitempty"`
	ExpiryTime time.Time `json:"expiry_time,omitempty"`
}

type SessionResult struct {
	Session Session
	Err     error
}

type AuthStatus int

const (
	NonAuthorised AuthStatus = iota
	Authorised
)

type AuthStatusResult struct {
	Status  AuthStatus `json:"status,omitempty"`
	Session Session    `json:"session,omitempty"`
	Err     error      `json:"error,omitempty"`
}

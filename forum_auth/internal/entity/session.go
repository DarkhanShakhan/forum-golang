package entity

type Session struct {
	Token      string `json:"token"`
	UserId     int64  `json:"user_id"`
	ExpiryDate string `json:"expiry_date"`
}

type SessionResult struct {
	Session Session
	Err     error
}

type AuthStatus int

const (
	Authorised AuthStatus = iota
	NonAuthorised
)

type AuthStatusResult struct {
	Status AuthStatus
	Err    error
}

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
	NonAuthorised AuthStatus = iota
	Authorised
)

type AuthStatusResult struct {
	Status AuthStatus `json:"status,omitempty"`
	Token  string     `json:"token,omitempty"`
	Err    error      `json:"error,omitempty"`
}

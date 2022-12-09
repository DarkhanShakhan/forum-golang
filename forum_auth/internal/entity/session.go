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

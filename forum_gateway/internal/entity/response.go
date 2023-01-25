package entity

type Response struct {
	Err          error
	UserId       int64       `json:"user_id,omitempty"`
	ErrorMessage string      `json:"error,omitempty"`
	AuthStatus   bool        `json:"authorised,omitempty"`
	Body         interface{} `json:"body,omitempty"`
}

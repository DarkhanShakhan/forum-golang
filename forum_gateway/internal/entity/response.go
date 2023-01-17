package entity

type Response struct {
	Err          error
	ErrorMessage string      `json:"error,omitempty"`
	AuthStatus   bool        `json:"authorised,omitempty"`
	Body         interface{} `json:"body,omitempty"`
}

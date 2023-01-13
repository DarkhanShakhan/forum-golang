package entity

type Response struct {
	ErrorMessage string      `json:"error,omitempty"`
	Body         interface{} `json:"body,omitempty"`
}

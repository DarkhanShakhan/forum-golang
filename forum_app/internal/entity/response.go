package entity

type Response struct {
	ErrorMessage string      `json:"error,omitempty"`
	Content      interface{} `json:"content,omitempty"`
}

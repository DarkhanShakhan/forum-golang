package entity

type Response struct {
	Err     Error       `json:"error,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

type Error struct {
	StatusCode   int    `json:"code,omitempty"`
	ErrorMessage string `json:"message,omitempty"`
}

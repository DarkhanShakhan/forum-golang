package entity

type Response struct {
	Err     string      `json:"error,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

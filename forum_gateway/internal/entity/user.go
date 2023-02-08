package entity

type User struct {
	Id int64 `json:"id,omitempty"`
}

type Result struct {
	Id  int
	Err error
}

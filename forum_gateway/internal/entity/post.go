package entity

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Post struct {
	Id       int        `json:"id,omitempty"`
	User     User       `json:"user,omitempty"`
	Title    string     `json:"title,omitempty"`
	Content  string     `json:"content,omitempty"`
	Category []Category `json:"categories,omitempty"`
}

type Category struct {
	Id int `json:"id,omitempty"`
}

type PostResult struct {
	Post Post
	Err  error
}

func GetPost(r *http.Request) (Post, error) {
	post := Post{}
	post.Title = r.FormValue("title")
	if strings.TrimSpace(post.Title) == "" {
		return Post{}, errors.New("Empty title")
	}
	cats := r.Form["category"]
	if len(cats) == 0 {
		return Post{}, errors.New("No category has been chosen")
	}
	for _, cat := range cats {
		cat_id, err := strconv.Atoi(cat)
		if err != nil {
			return Post{}, errors.New("Invalid category")
		}
		post.Category = append(post.Category, Category{Id: cat_id})
	}
	post.Content = r.FormValue("content")
	return post, nil
}

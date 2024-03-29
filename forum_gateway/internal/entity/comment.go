package entity

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Comment struct {
	Id      int    `json:"id,omitempty"`
	Post    Post   `json:"post,omitempty"`
	User    User   `json:"user,omitempty"`
	Content string `json:"comment_content,omitempty"`
}

func GetComment(r *http.Request) CommentResult {
	commentRes := CommentResult{}
	var user_id interface{} = r.Context().Value("user_id")
	commentRes.Comment.User.Id = user_id.(int64)
	if commentRes.Comment.User.Id == 0 {
		return CommentResult{Err: errors.New("User is not provided")}
	}
	post_id, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		return CommentResult{Err: err}
	}
	commentRes.Comment.Post = Post{Id: post_id}
	commentRes.Comment.Content = r.FormValue("content")
	if strings.TrimSpace(commentRes.Comment.Content) == "" {
		return CommentResult{Err: ErrEmptyComment}
	}
	return commentRes
}

type CommentResult struct {
	Comment Comment
	Err     error
}

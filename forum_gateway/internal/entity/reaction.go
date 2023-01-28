package entity

import (
	"errors"
	"net/http"
	"strconv"
)

type Reaction struct {
	Like bool   `json:"like,omitempty"`
	Date string `json:"reaction_date,omitempty"`
	User User   `json:"user,omitempty"`
}

type PostReaction struct {
	Reaction `json:"reaction,omitempty"`
	Post     `json:"post,omitempty"`
}

func GetPostReaction(r *http.Request) (PostReaction, error) {
	var (
		postReaction PostReaction
		err          error
		id           interface{} = r.Context().Value("user_id")
		ok           bool
		reaction     string = r.FormValue("reaction")
		post_id      string = r.FormValue("post_id")
	)
	postReaction.Post.Id, err = strconv.Atoi(post_id)
	if err != nil {
		return PostReaction{}, err
	}
	postReaction.Reaction.User.Id, ok = id.(int64)
	if !ok {
		return PostReaction{}, errors.New("invalid user id")
	}
	switch reaction {
	case "true":
		postReaction.Reaction.Like = true
	case "false":
		postReaction.Reaction.Like = false
	default:
		return PostReaction{}, errors.New("invalid reaction")
	}
	return postReaction, nil
}

func GetCommentReaction(r *http.Request) (CommentReaction, error) {
	var (
		commentReaction CommentReaction
		err             error
		id              interface{} = r.Context().Value("user_id")
		ok              bool
		reaction        string = r.FormValue("reaction")
		post_id         string = r.FormValue("post_id")
		comment_id      string = r.FormValue("comment_id")
	)
	commentReaction.Id, err = strconv.Atoi(comment_id)
	if err != nil {
		return CommentReaction{}, err
	}
	commentReaction.Post.Id, err = strconv.Atoi(post_id)
	if err != nil {
		return CommentReaction{}, err
	}
	commentReaction.Reaction.User.Id, ok = id.(int64)
	if !ok {
		return CommentReaction{}, errors.New("invalid user id")
	}
	switch reaction {
	case "true":
		commentReaction.Reaction.Like = true
	case "false":
		commentReaction.Reaction.Like = false
	default:
		return CommentReaction{}, errors.New("invalid reaction")
	}
	return commentReaction, nil
}

type CommentReaction struct {
	Reaction `json:"reaction,omitempty"`
	Comment  `json:"comment,omitempty"`
}

type ReactionsResult struct {
	Reactions []Reaction
	Err       error
}

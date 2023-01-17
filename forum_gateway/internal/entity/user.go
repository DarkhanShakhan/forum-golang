package entity

type User struct {
	Id                   int64             `json:"id,omitempty"`
	Name                 string            `json:"name,omitempty"`
	Email                string            `json:"email,omitempty"`
	Password             string            `json:"password,omitempty"`
	RegDate              string            `json:"registration_date,omitempty"`
	Posts                []Post            `json:"posts,omitempty"`
	TotalPosts           int               `json:"total_posts,omitempty"`
	Comments             []Comment         `json:"comments,omitempty"`
	TotalComments        int               `json:"total_comments,omitempty"`
	PostLikes            []PostReaction    `json:"post_likes,omitempty"`
	TotalPostLikes       int               `json:"total_post_reactions,omitempty"`
	PostDislikes         []PostReaction    `json:"post_dislikes,omitempty"`
	TotalPostDislikes    int               `json:"total_post_dislikes,omitempty"`
	CommentLikes         []CommentReaction `json:"comment_likes,omitempty"`
	TotalCommentLikes    int               `json:"total_comment_likes,omitempty"`
	CommentDislikes      []CommentReaction `json:"comment_dislikes,omitempty"`
	TotalCommentDislikes int               `json:"total_comment_dislikes,omitempty"`
}

func (u *User) CountTotals() {
	u.TotalPosts = len(u.Posts)
	u.TotalComments = len(u.Comments)
	u.TotalPostLikes = len(u.PostLikes)
	u.TotalPostDislikes = len(u.PostDislikes)
	u.TotalCommentLikes = len(u.CommentLikes)
	u.TotalCommentDislikes = len(u.CommentDislikes)
}

type Result struct {
	Id  int
	Err error
}

type UserResult struct {
	User User
	Err  error
}

type UsersResult struct {
	Users []User
	Err   error
}

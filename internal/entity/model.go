package entity

type Post struct {
	Id            int        `json:"id,omitempty"`
	User          User       `json:"user,omitempty"`
	Date          string     `json:"date,omitempty"`
	Title         string     `json:"title,omitempty"`
	Content       string     `json:"content,omitempty"`
	Category      []Category `json:"categories,omitempty"`
	Comments      []Comment  `json:"comments,omitempty"`
	CountComments int        `json:"total_comments,omitempty"`
	Likes         []Reaction `json:"likes,omitempty"`
	CountLikes    int        `json:"total_likes,omitempty"`
	Dislikes      []Reaction `json:"dislkes,omitempty"`
	CountDislikes int        `json:"total_dislikes,omitempty"`
}

func (p *Post) CountTotals() {
	p.CountComments = len(p.Comments)
	p.CountLikes = len(p.Likes)
	p.CountDislikes = len(p.Dislikes)
}

type Comment struct {
	Id            int        `json:"id,omitempty"`
	Post          Post       `json:"post,omitempty"`
	User          User       `json:"user,omitempty"`
	Date          string     `json:"comment_date,omitempty"`
	Content       string     `json:"comment_content,omitempty"`
	Likes         []Reaction `json:"likes,omitempty"`
	TotalLikes    int        `json:"total_likes,omitempty"`
	Dislikes      []Reaction `json:"dislikes,omitempty"`
	TotalDislikes int        `json:"total_dislikes,omitempty"`
}

func (c *Comment) CountTotals() {
	c.TotalLikes = len(c.Likes)
	c.TotalDislikes = len(c.Dislikes)
}

type Reaction struct {
	Like bool   `json:"like,omitempty"`
	Date string `json:"reaction_date,omitempty"`
	User User   `json:"user,omitempty"`
}

type PostReaction struct {
	Reaction `json:"reaction,omitempty"`
	Post     `json:"post,omitempty"`
}

type CommentReaction struct {
	Reaction `json:"reaction,omitempty"`
	Comment  `json:"comment,omitempty"`
}

type User struct {
	Id                   int               `json:"id,omitempty"`
	Name                 string            `json:"name,omitempty"`
	Email                string            `json:"email,omitempty"`
	Password             string            `json:"password,omitempty"`
	RegDate              string            `json:"registration_date,omitempty"`
	DateOfBirth          string            `json:"date_of_birth,omitempty"`
	City                 string            `json:"city,omitempty"`
	Sex                  string            `json:"sex,omitempty"`
	Posts                []Post            `json:"posts,omitempty"`
	CountPosts           int               `json:"total_posts,omitempty"`
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
	u.CountPosts = len(u.Posts)
	u.TotalComments = len(u.Comments)
	u.TotalPostLikes = len(u.PostLikes)
	u.TotalPostDislikes = len(u.PostDislikes)
	u.TotalCommentLikes = len(u.CommentLikes)
	u.TotalCommentDislikes = len(u.CommentDislikes)
}

type Category struct {
	Title      string `json:"title,omitempty"`
	Posts      []Post `json:"posts,omitempty"`
	TotalPosts int    `json:"total_posts,omitempty"`
}

func (c *Category) CountTotals() {
	c.TotalPosts = len(c.Posts)
}

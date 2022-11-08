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
	CountLikes    int        `json:"total_likes,omitempty"`
	Dislikes      []Reaction `json:"dislikes,omitempty"`
	CountDislikes int        `json:"total_dislikes,omitempty"`
}

func (c *Comment) CountTotals() {
	c.CountLikes = len(c.Likes)
	c.CountDislikes = len(c.Dislikes)
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
	Id                    int               `json:"id,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Email                 string            `json:"email,omitempty"`
	Password              string            `json:"password,omitempty"`
	RegDate               string            `json:"registration_date,omitempty"`
	DateOfBirth           string            `json:"date_of_birth,omitempty"`
	City                  string            `json:"city,omitempty"`
	Sex                   string            `json:"sex,omitempty"`
	Posts                 []Post            `json:"posts,omitempty"`
	CountPosts            int               `json:"total_posts,omitempty"`
	Comments              []Comment         `json:"comments,omitempty"`
	TotalComments         int               `json:"total_comments,omitempty"`
	PostReactions         []PostReaction    `json:"post_reactions,omitempty"`
	CountPostReactions    int               `json:"total_post_reactions,omitempty"`
	CommentReactions      []CommentReaction `json:"comment_reactions,omitempty"`
	CountCommentReactions int               `json:"total_comment_reactions,omitempty"`
}

func (u *User) CountTotals() {
	u.CountPosts = len(u.Posts)
	u.TotalComments = len(u.Comments)
	u.CountPostReactions = len(u.PostReactions)
	u.CountCommentReactions = len(u.CommentReactions)
}

type Category struct {
	Title      string `json:"title,omitempty"`
	Posts      []Post `json:"posts,omitempty"`
	CountPosts int    `json:"total_posts,omitempty"`
}

func (c *Category) CountTotals() {
	c.CountPosts = len(c.Posts)
}

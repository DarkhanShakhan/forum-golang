package entity

type Post struct {
	Id            int        `json:"id,omitempty"`
	User          User       `json:"user,omitempty"`
	Date          string     `json:"date,omitempty"`
	Title         string     `json:"title,omitempty"`
	Content       string     `json:"content,omitempty"`
	Category      []Category `json:"categories,omitempty"`
	Comments      []Comment  `json:"comments,omitempty"`
	TotalComments int        `json:"total_comments,omitempty"`
	Likes         []Reaction `json:"likes,omitempty"`
	TotalLikes    int        `json:"total_likes,omitempty"`
	Dislikes      []Reaction `json:"dislkes,omitempty"`
	TotalDislikes int        `json:"total_dislikes,omitempty"`
}

func (p *Post) CountTotals() {
	p.TotalComments = len(p.Comments)
	p.TotalLikes = len(p.Likes)
	p.TotalDislikes = len(p.Dislikes)
}

type PostResult struct {
	Post Post
	Err  error
}

type PostsResult struct {
	Posts []Post
	Err   error
}

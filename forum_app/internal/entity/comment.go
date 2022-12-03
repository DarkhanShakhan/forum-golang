package entity

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

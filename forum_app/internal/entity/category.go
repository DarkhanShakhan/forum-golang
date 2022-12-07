package entity

type Category struct {
	Id         int    `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Posts      []Post `json:"posts,omitempty"`
	TotalPosts int    `json:"total_posts,omitempty"`
}

func (c *Category) CountTotals() {
	c.TotalPosts = len(c.Posts)
}

type CatResult struct {
	Cat Category
	Err error
}

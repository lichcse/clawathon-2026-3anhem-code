package entity

// PostStatusDefault new post
const PostStatusDefault int = 0
const PostStatusPublished int = 1
const PostStatusDraft int = 2

// Post struct post entity
type Post struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TableName func get table name
func (p *Post) TableName() string {
	return "im_post"
}
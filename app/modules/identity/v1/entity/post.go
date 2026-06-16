package entity

const (
	// PostStatusDraft post is draft
	PostStatusDraft int = 0
	// PostStatusPublished post is published
	PostStatusPublished int = 1
	// PostStatusDeleted post is deleted
	PostStatusDeleted int = 2
)

// Post struct post entity
type Post struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	ImageURL    string `json:"image_url"`
	AuthorID    uint64 `json:"author_id"`
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
	Status      int    `json:"status"`
}

// TableName func get table name
func (p *Post) TableName() string {
	return "im_post"
}

package entity

const (
	// PostStatusDraft post is draft
	PostStatusDraft int = 0
	// PostStatusPublished post is published
	PostStatusPublished int = 1
	// PostStatusArchived post is archived
	PostStatusArchived int = 2
)

// Post struct post entity
type Post struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"user_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	ImageURL    string `json:"image_url"`
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
	PublishedAt string `json:"published_at"`
	Status      int    `json:"status"`
}

// TableName func get table name
func (p *Post) TableName() string {
	return "im_post"
}

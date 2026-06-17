package v1

// PostAddRequest post add request data
type PostAddRequest struct {
	Title    string `json:"title" swaggertype:"string" maxLength:"3" maxLength:"200" example:"My First Post"`
	Content  string `json:"content" swaggertype:"string" example:"This is the content of my post"`
	ImageURL string `json:"image_url" swaggertype:"string" example:"https://example.com/image.jpg"`
}

// PostUpdateRequest post update request data
type PostUpdateRequest struct {
	Title    string `json:"title" swaggertype:"string" maxLength:"3" maxLength:"200" example:"My Updated Post"`
	Content  string `json:"content" swaggertype:"string" example:"This is the updated content"`
	ImageURL string `json:"image_url" swaggertype:"string" example:"https://example.com/new-image.jpg"`
	Status   int    `json:"status" swaggertype:"integer" example:"1"`
}

// PostDetailResponse post response data
type PostDetailResponse struct {
	ID          uint64 `json:"id" swaggertype:"integer" example:"1"`
	UserID      uint64 `json:"user_id" swaggertype:"integer" example:"1"`
	Title       string `json:"title" swaggertype:"string" example:"My First Post"`
	Content     string `json:"content" swaggertype:"string" example:"This is the content of my post"`
	ImageURL    string `json:"image_url" swaggertype:"string" example:"https://example.com/image.jpg"`
	CreatedAt   string `json:"created_at" swaggertype:"string" example:"2024-01-01 10:00:00"`
	ModifiedAt  string `json:"modified_at" swaggertype:"string" example:"2024-01-01 10:00:00"`
	PublishedAt string `json:"published_at" swaggertype:"string" example:"2024-01-01 10:00:00"`
	Status      int    `json:"status" swaggertype:"integer" example:"1"`
}

// PostListRequest post list request data
type PostListRequest struct {
	Page  int `json:"page" form:"page" swaggertype:"integer" example:"1"`
	Limit int `json:"limit" form:"limit" swaggertype:"integer" example:"10"`
}

// PostListResponse post list response data
type PostListResponse struct {
	Posts     []PostDetailResponse `json:"posts"`
	Total     int64                `json:"total" swaggertype:"integer" example:"100"`
	Page      int                  `json:"page" swaggertype:"integer" example:"1"`
	Limit     int                  `json:"limit" swaggertype:"integer" example:"10"`
	TotalPage int                  `json:"total_page" swaggertype:"integer" example:"10"`
}

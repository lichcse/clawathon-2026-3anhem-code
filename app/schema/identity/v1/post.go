package v1

// PostAddRequest post add request data
type PostAddRequest struct {
	Title    string `json:"title" swaggertype:"string" example:"My First Post"`
	Content  string `json:"content" swaggertype:"string" example:"This is the content of my post"`
	ImageURL string `json:"image_url" swaggertype:"string" example:"https://example.com/image.jpg"`
}

// PostUpdateRequest post update request data
type PostUpdateRequest struct {
	Title    string `json:"title" swaggertype:"string" example:"Updated Title"`
	Content  string `json:"content" swaggertype:"string" example:"Updated content"`
	ImageURL string `json:"image_url" swaggertype:"string" example:"https://example.com/new-image.jpg"`
}

// PostDetailResponse post response data
type PostDetailResponse struct {
	ID          uint64 `json:"id" swaggertype:"integer" example:"1"`
	Title       string `json:"title" swaggertype:"string" example:"My First Post"`
	Content     string `json:"content" swaggertype:"string" example:"This is the content"`
	ImageURL    string `json:"image_url" swaggertype:"string" example:"https://example.com/image.jpg"`
	AuthorID    uint64 `json:"author_id" swaggertype:"integer" example:"1"`
	CreatedAt   string `json:"created_at" swaggertype:"string" example:"1991-02-13 10:10:10"`
	ModifiedAt  string `json:"modified_at" swaggertype:"string" example:"2020-07-15 10:10:10"`
	Status      int    `json:"status" swaggertype:"integer" example:"1"`
}

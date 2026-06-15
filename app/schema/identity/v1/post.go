package v1

// PostAddRequest post add request data
type PostAddRequest struct {
	Title   string `json:"title" swaggertype:"string" example:"Post Title"`
	Content string `json:"content" swaggertype:"string" example:"Post content"`
}

// PostUpdateRequest post update request data
type PostUpdateRequest struct {
	Title   string `json:"title" swaggertype:"string" example:"Post Title"`
	Content string `json:"content" swaggertype:"string" example:"Post content"`
	Status  int    `json:"status" swaggertype:"integer" example:"1"`
}

// PostDetailResponse post response data
type PostDetailResponse struct {
	ID        uint64 `json:"id" swaggertype:"integer" example:"1"`
	UserID    uint64 `json:"user_id" swaggertype:"integer" example:"1"`
	Title     string `json:"title" swaggertype:"string" example:"Post Title"`
	Content   string `json:"content" swaggertype:"string" example:"Post content"`
	Status    int    `json:"status" swaggertype:"integer" example:"1"`
	CreatedAt string `json:"created_at" swaggertype:"string" example:"1991-02-13 10:10:10"`
	UpdatedAt string `json:"updated_at" swaggertype:"string" example:"2020-07-15 10:10:10"`
}
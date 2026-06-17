package service

import (
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/entity"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/repository"
	schema "github.com/lichcse/clawathon-2026-3anhem-code/app/schema/post/v1"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"
)

// PostService interface of post service object
type PostService interface {
	Add(userID uint64, postAddRequest *schema.PostAddRequest) (*schema.PostDetailResponse, error)
	Detail(id uint64) (*schema.PostDetailResponse, error)
	Update(id string, postUpdate *schema.PostUpdateRequest) error
	Delete(id string) error
	List(userID uint64, page, limit int) (*schema.PostListResponse, error)
}

type postService struct {
	postRepo repository.PostRepository
	convert  utils.IMConvert
}

// NewPostService func new post service object
func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{
		postRepo: postRepo,
		convert:  utils.NewIMConvert(),
	}
}

// Add func add new post
func (s *postService) Add(userID uint64, postAddRequest *schema.PostAddRequest) (*schema.PostDetailResponse, error) {
	result := &schema.PostDetailResponse{}
	post := &entity.Post{
		UserID:   userID,
		Title:    postAddRequest.Title,
		Content:  postAddRequest.Content,
		ImageURL: postAddRequest.ImageURL,
	}

	err := s.postRepo.Add(post)
	if err != nil {
		return result, err
	}

	err = s.convert.Object(post, &result)
	return result, err
}

// Detail func get detail post info
func (s *postService) Detail(id uint64) (*schema.PostDetailResponse, error) {
	result := &schema.PostDetailResponse{}
	post, err := s.postRepo.Detail(id)
	if err != nil {
		return result, s.convert.DatabaseError(err)
	}

	err = s.convert.Object(post, &result)
	return result, err
}

// Update func update post info
func (s *postService) Update(id string, postUpdate *schema.PostUpdateRequest) error {
	post := &entity.Post{
		Title:    postUpdate.Title,
		Content:  postUpdate.Content,
		ImageURL: postUpdate.ImageURL,
		Status:   postUpdate.Status,
	}

	// Set published_at if status is published
	if postUpdate.Status == entity.PostStatusPublished {
		post.PublishedAt = utils.NewIMTime().TimeDB()
	}

	return s.postRepo.Update(id, post)
}

// Delete func delete post info
func (s *postService) Delete(id string) error {
	return s.postRepo.Delete(id)
}

// List func get list post with pagination
func (s *postService) List(userID uint64, page, limit int) (*schema.PostListResponse, error) {
	result := &schema.PostListResponse{}

	// Set default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	posts, total, err := s.postRepo.List(userID, page, limit)
	if err != nil {
		return result, s.convert.DatabaseError(err)
	}

	// Convert posts to response
	var postResponses []schema.PostDetailResponse
	err = s.convert.Object(posts, &postResponses)
	if err != nil {
		return result, err
	}

	// Calculate total page
	totalPage := int(total) / limit
	if int(total)%limit > 0 {
		totalPage++
	}

	result.Posts = postResponses
	result.Total = total
	result.Page = page
	result.Limit = limit
	result.TotalPage = totalPage

	return result, nil
}

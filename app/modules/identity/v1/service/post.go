package service

import (
	"errors"
	"strconv"

	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/entity"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/repository"
	schema "github.com/lichcse/clawathon-2026-3anhem-code/app/schema/identity/v1"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"
)

// PostService interface of post service object
type PostService interface {
	Add(userID uint64, postAddRequest *schema.PostAddRequest) (*schema.PostDetailResponse, error)
	Detail(id uint64) (*schema.PostDetailResponse, error)
	Update(id string, postUpdate *schema.PostUpdateRequest) error
	Delete(id string) error
	ListByUserID(userID uint64) ([]*schema.PostDetailResponse, error)
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
func (p *postService) Add(userID uint64, postAddRequest *schema.PostAddRequest) (*schema.PostDetailResponse, error) {
	result := &schema.PostDetailResponse{}
	post := &entity.Post{
		UserID:  userID,
		Title:   postAddRequest.Title,
		Content: postAddRequest.Content,
		Status:  entity.PostStatusDefault,
	}

	err := p.postRepo.Add(post)
	if err != nil {
		return result, err
	}

	err = p.convert.Object(post, &result)
	return result, err
}

// Detail func get detail post info
func (p *postService) Detail(id uint64) (*schema.PostDetailResponse, error) {
	result := &schema.PostDetailResponse{}
	post, err := p.postRepo.Detail(id)
	if err != nil {
		return result, p.convert.DatabaseError(err)
	}

	err = p.convert.Object(post, &result)
	return result, err
}

// Update func update post info
func (p *postService) Update(id string, postUpdate *schema.PostUpdateRequest) error {
	post := &entity.Post{
		Title:   postUpdate.Title,
		Content: postUpdate.Content,
		Status:  postUpdate.Status,
	}
	return p.postRepo.Update(id, post)
}

// Delete func delete post info
func (p *postService) Delete(id string) error {
	// Check if post exists before deleting
	postID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errors.New("post_not_found")
	}

	post, err := p.postRepo.Detail(postID)
	if err != nil {
		return errors.New("post_not_found")
	}

	if post == nil || post.ID == 0 {
		return errors.New("post_not_found")
	}

	return p.postRepo.Delete(id)
}

// ListByUserID func get list post by user id
func (p *postService) ListByUserID(userID uint64) ([]*schema.PostDetailResponse, error) {
	posts, err := p.postRepo.ListByUserID(userID)
	if err != nil {
		return nil, p.convert.DatabaseError(err)
	}

	var result []*schema.PostDetailResponse
	for _, post := range posts {
		var item schema.PostDetailResponse
		p.convert.Object(post, &item)
		result = append(result, &item)
	}
	return result, nil
}
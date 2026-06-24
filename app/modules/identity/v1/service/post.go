package service

import (
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/entity"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/repository"
	schema "github.com/lichcse/clawathon-2026-3anhem-code/app/schema/identity/v1"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"
)

// PostService interface of post service object
type PostService interface {
	Add(postAddRequest *schema.PostAddRequest, authorID uint64) (*schema.PostDetailResponse, error)
	Detail(id uint64) (*schema.PostDetailResponse, error)
	Update(id string, postUpdate *schema.PostUpdateRequest) error
	Delete(id string) error
	ListByAuthor(authorID uint64, limit, offset int) ([]schema.PostDetailResponse, error)
	Publish(id string) error
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
func (p *postService) Add(postAddRequest *schema.PostAddRequest, authorID uint64) (*schema.PostDetailResponse, error) {
	result := &schema.PostDetailResponse{}
	post := &entity.Post{
		Title:    postAddRequest.Title,
		Content:  postAddRequest.Content,
		ImageURL: postAddRequest.ImageURL,
		AuthorID: authorID,
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
		Title:    postUpdate.Title,
		Content:  postUpdate.Content,
		ImageURL: postUpdate.ImageURL,
	}
	return p.postRepo.Update(id, post)
}

// Delete func delete post info
func (p *postService) Delete(id string) error {
	return p.postRepo.Delete(id)
}

// ListByAuthor func get list post by author
func (p *postService) ListByAuthor(authorID uint64, limit, offset int) ([]schema.PostDetailResponse, error) {
	posts, err := p.postRepo.ListByAuthor(authorID, limit, offset)
	if err != nil {
		return nil, p.convert.DatabaseError(err)
	}

	var result []schema.PostDetailResponse
	for _, post := range posts {
		var item schema.PostDetailResponse
		p.convert.Object(post, &item)
		result = append(result, item)
	}
	return result, nil
}

// Publish func publish a post
func (p *postService) Publish(id string) error {
	post := &entity.Post{
		Status: entity.PostStatusPublished,
	}
	return p.postRepo.Update(id, post)
}

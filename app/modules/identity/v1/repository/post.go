package repository

import (
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/entity"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"

	"gorm.io/gorm"
)

// PostRepository interface of post repository object
type PostRepository interface {
	Add(post *entity.Post) error
	Detail(id uint64) (*entity.Post, error)
	Update(id string, post *entity.Post) error
	Delete(id string) error
	ListByUserID(userID uint64) ([]*entity.Post, error)
}

type postRepository struct {
	db     *gorm.DB
	imTime utils.IMTime
}

// NewPostRepository func new post repository object
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db, imTime: utils.NewIMTime()}
}

// Add func add new post
func (p *postRepository) Add(post *entity.Post) error {
	post.Status = entity.PostStatusDefault
	post.CreatedAt = p.imTime.TimeDB()
	post.UpdatedAt = p.imTime.TimeDB()
	return p.db.Create(&post).Error
}

// Detail func get detail post info
func (p *postRepository) Detail(id uint64) (*entity.Post, error) {
	result := &entity.Post{}
	err := p.db.Where("id = ?", id).First(result).Error
	return result, err
}

// Update func update post info
func (p *postRepository) Update(id string, post *entity.Post) error {
	post.UpdatedAt = p.imTime.TimeDB()
	return p.db.Model(&entity.Post{}).Where("id = ?", id).Updates(post).Error
}

// Delete func delete post info
func (p *postRepository) Delete(id string) error {
	return p.db.Where("id = ?", id).Delete(&entity.Post{}).Error
}

// ListByUserID func get list post by user id
func (p *postRepository) ListByUserID(userID uint64) ([]*entity.Post, error) {
	var posts []*entity.Post
	err := p.db.Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}
package repository

import (
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/entity"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"

	"gorm.io/gorm"
)

// PostRepository interface of post repository object
type PostRepository interface {
	Add(post *entity.Post) error
	Detail(id uint64) (*entity.Post, error)
	Update(id string, post *entity.Post) error
	Delete(id string) error
	List(userID uint64, page, limit int) ([]entity.Post, int64, error)
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
func (r *postRepository) Add(post *entity.Post) error {
	post.Status = entity.PostStatusDraft
	post.CreatedAt = r.imTime.TimeDB()
	post.ModifiedAt = r.imTime.TimeDB()
	return r.db.Create(&post).Error
}

// Detail func get detail post info
func (r *postRepository) Detail(id uint64) (*entity.Post, error) {
	result := &entity.Post{}
	err := r.db.Where("id = ?", id).First(result).Error
	return result, err
}

// Update func update post info
func (r *postRepository) Update(id string, post *entity.Post) error {
	post.ModifiedAt = r.imTime.TimeDB()
	return r.db.Where("id = ?", id).Updates(post).Error
}

// Delete func delete post info
func (r *postRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&entity.Post{}).Error
}

// List func get list post with pagination
func (r *postRepository) List(userID uint64, page, limit int) ([]entity.Post, int64, error) {
	var posts []entity.Post
	var total int64

	query := r.db.Model(&entity.Post{})

	// Filter by userID if provided
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get posts with pagination
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	return posts, total, err
}

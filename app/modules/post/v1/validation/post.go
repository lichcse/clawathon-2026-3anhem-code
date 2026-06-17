package validation

import (
	"errors"

	schema "github.com/lichcse/clawathon-2026-3anhem-code/app/schema/post/v1"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"
)

// PostValidation interface of post validation object
type PostValidation interface {
	Add(postAdd *schema.PostAddRequest) error
	Update(postUpdate *schema.PostUpdateRequest) error
}

type postValidation struct {
	validation utils.IMValidation
}

// NewPostValidation func new post validation object
func NewPostValidation() PostValidation {
	return &postValidation{validation: utils.NewIMValidation()}
}

// Add func validate data add
func (v *postValidation) Add(postAddRequest *schema.PostAddRequest) error {
	if len(postAddRequest.Title) < 3 || len(postAddRequest.Title) > 200 {
		return errors.New("post_invalid_title")
	}

	if len(postAddRequest.Content) == 0 {
		return errors.New("post_invalid_content")
	}

	return nil
}

// Update func validate data update
func (v *postValidation) Update(postUpdateRequest *schema.PostUpdateRequest) error {
	if len(postUpdateRequest.Title) > 0 && (len(postUpdateRequest.Title) < 3 || len(postUpdateRequest.Title) > 200) {
		return errors.New("post_invalid_title")
	}

	// Validate status if provided
	if postUpdateRequest.Status != 0 {
		if postUpdateRequest.Status != 1 && postUpdateRequest.Status != 2 {
			return errors.New("post_invalid_status")
		}
	}

	return nil
}

package validation

import (
	"errors"

	schema "github.com/lichcse/clawathon-2026-3anhem-code/app/schema/identity/v1"
)

// PostValidation interface of post validation object
type PostValidation interface {
	Add(postAdd *schema.PostAddRequest) error
	Update(postUpdate *schema.PostUpdateRequest) error
}

type postValidation struct {
}

// NewPostValidation func new post validation object
func NewPostValidation() PostValidation {
	return &postValidation{}
}

// Add func validate data add
func (p *postValidation) Add(postAddRequest *schema.PostAddRequest) error {
	if len(postAddRequest.Title) == 0 {
		return errors.New("post_invalid_title")
	}

	if len(postAddRequest.Title) > 200 {
		return errors.New("post_title_too_long")
	}

	if len(postAddRequest.Content) == 0 {
		return errors.New("post_invalid_content")
	}

	return nil
}

// Update func validate data update
func (p *postValidation) Update(postUpdateRequest *schema.PostUpdateRequest) error {
	if len(postUpdateRequest.Title) == 0 {
		return errors.New("post_invalid_title")
	}

	if len(postUpdateRequest.Title) > 200 {
		return errors.New("post_title_too_long")
	}

	if len(postUpdateRequest.Content) == 0 {
		return errors.New("post_invalid_content")
	}

	return nil
}

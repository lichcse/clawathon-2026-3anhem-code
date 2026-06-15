package validation

import (
	"errors"

	schema "github.com/lichcse/clawathon-2026-3anhem-code/app/schema/identity/v1"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"
)

// UserValidation interface of user validation object
type UserValidation interface {
	Add(userAdd *schema.UserAddRequest) error
}

type userValidation struct {
	validation utils.IMValidation
}

// NewUserValidation func new user validation object
func NewUserValidation() UserValidation {
	return &userValidation{validation: utils.NewIMValidation()}
}

// Add func validate data add
func (u *userValidation) Add(userAddRequest *schema.UserAddRequest) error {

	if len(userAddRequest.FullName) < 3 || len(userAddRequest.FullName) > 50 {
		return errors.New("user_invalid_full_name")
	}

	validUsername := u.validation.IsValidUsername(userAddRequest.Username)
	if !validUsername {
		return errors.New("user_invalid_username")
	}

	validEmail := u.validation.IsValidEmail(userAddRequest.Email)
	if !validEmail {
		return errors.New("user_invalid_email")
	}

	if len(userAddRequest.Password) == 0 {
		return errors.New("user_invalid_password")
	}
	return nil
}

package resource

import (
	"net/http"

	"github.com/lichcse/clawathon-2026-3anhem-code/app/response"
)

// CodeMessageMapping data
// code format: xxx.yyy.zzz
//   - xxx: module
//   - yyy: entity
//   - zzz: entity error code
//
// module:
//   - identity: 001
//
// entity:
//   - general: 000
//   - user: 001
//   - post: 002
var CodeMessageMapping = response.CodeMessageMapping{
	"success": response.CodeStatus{
		Code:   "001.000.000",
		Status: http.StatusOK,
	},
	"not_allow": response.CodeStatus{
		Code:   "001.000.001",
		Status: http.StatusBadRequest,
	},
	"user_invalid_full_name": response.CodeStatus{
		Code:   "001.001.001",
		Status: http.StatusBadRequest,
	},
	"user_invalid_username": response.CodeStatus{
		Code:   "001.001.002",
		Status: http.StatusBadRequest,
	},
	"user_invalid_email": response.CodeStatus{
		Code:   "001.001.003",
		Status: http.StatusBadRequest,
	},
	"user_invalid_password": response.CodeStatus{
		Code:   "001.001.004",
		Status: http.StatusBadRequest,
	},
	"post_invalid_title": response.CodeStatus{
		Code:   "001.002.001",
		Status: http.StatusBadRequest,
	},
	"post_title_too_long": response.CodeStatus{
		Code:   "001.002.002",
		Status: http.StatusBadRequest,
	},
	"post_invalid_content": response.CodeStatus{
		Code:   "001.002.003",
		Status: http.StatusBadRequest,
	},
	"post_not_found": response.CodeStatus{
		Code:   "001.002.004",
		Status: http.StatusNotFound,
	},
}
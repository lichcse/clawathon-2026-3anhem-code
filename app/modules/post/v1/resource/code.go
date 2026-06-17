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
//   - post: 002
//
// entity:
//   - general: 000
//   - post: 001
var CodeMessageMapping = response.CodeMessageMapping{
	"success": response.CodeStatus{
		Code:   "002.000.000",
		Status: http.StatusOK,
	},
	"not_allow": response.CodeStatus{
		Code:   "002.000.001",
		Status: http.StatusBadRequest,
	},
	"post_invalid_title": response.CodeStatus{
		Code:   "002.001.001",
		Status: http.StatusBadRequest,
	},
	"post_invalid_content": response.CodeStatus{
		Code:   "002.001.002",
		Status: http.StatusBadRequest,
	},
	"post_invalid_status": response.CodeStatus{
		Code:   "002.001.003",
		Status: http.StatusBadRequest,
	},
}

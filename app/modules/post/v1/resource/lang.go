package resource

import (
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"
)

// DefaultLang default language
var DefaultLang = "en"

// PostLang post language
var PostLang = map[string]utils.Lang{
	"en": map[string]string{
		"success":                "Success",
		"not_allow":              "Not allowed",
		"post_invalid_title":     "Title must be between 3 and 200 characters",
		"post_invalid_content":   "Content is required",
		"post_invalid_status":    "Status must be 1 (published) or 2 (archived)",
	},
	"vi": map[string]string{
		"success":                "Thành công",
		"not_allow":              "Không được phép",
		"post_invalid_title":     "Tiêu đề phải từ 3 đến 200 ký tự",
		"post_invalid_content":   "Nội dung là bắt buộc",
		"post_invalid_status":    "Trạng thái phải là 1 (đã xuất bản) hoặc 2 (lưu trữ)",
	},
}

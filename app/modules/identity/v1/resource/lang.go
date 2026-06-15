package resource

import (
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/resource/lang"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"
)

// DefaultLang default language
var DefaultLang = "en"

// IdentityLang identity language
var IdentityLang = map[string]utils.Lang{
	"en": lang.EN,
	"vi": lang.VI,
}

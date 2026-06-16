package v1

import (
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/handler/http"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/resource"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/response"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"

	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/repository"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/service"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/validation"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var restResponse = response.NewRestResponse(utils.NewIMLanguage(resource.IdentityLang, resource.DefaultLang), resource.CodeMessageMapping)

// UserRoute : user route group
func UserRoute(r *gin.RouterGroup, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	userHandler := http.NewUserHandler(service.NewUserService(userRepository), restResponse, validation.NewUserValidation())
	r.POST("/user", userHandler.Add)
	r.GET("/user/:id", userHandler.Detail)
	r.PUT("/user/:id", userHandler.Update)
	r.DELETE("/user/:id", userHandler.Delete)
}

// PostRoute : post route group
func PostRoute(r *gin.RouterGroup, db *gorm.DB) {
	postRepository := repository.NewPostRepository(db)
	postHandler := http.NewPostHandler(service.NewPostService(postRepository), restResponse, validation.NewPostValidation())
	r.GET("/post", postHandler.List)
	r.POST("/post", postHandler.Add)
	r.GET("/post/:id", postHandler.Detail)
	r.PUT("/post/:id", postHandler.Update)
	r.DELETE("/post/:id", postHandler.Delete)
	r.POST("/post/:id/publish", postHandler.Publish)
}

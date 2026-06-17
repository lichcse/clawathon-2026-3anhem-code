package v1

import (
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/handler/http"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/repository"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/service"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/validation"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/response"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/utils"

	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/resource"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var postResponse = response.NewRestResponse(utils.NewIMLanguage(resource.PostLang, resource.DefaultLang), resource.CodeMessageMapping)

// PostRoute : post route group
func PostRoute(r *gin.RouterGroup, db *gorm.DB) {
	postRepository := repository.NewPostRepository(db)
	postHandler := http.NewPostHandler(
		service.NewPostService(postRepository),
		postResponse,
		validation.NewPostValidation(),
	)

	// Public routes
	r.GET("/post", postHandler.List)
	r.GET("/post/:id", postHandler.Detail)

	// Protected routes (require auth)
	protected := r.Group("")
	protected.Use(func(ctx *gin.Context) {
		// TODO: Add auth middleware here
		ctx.Next()
	})
	protected.POST("/post", postHandler.Add)
	protected.PUT("/post/:id", postHandler.Update)
	protected.DELETE("/post/:id", postHandler.Delete)
}

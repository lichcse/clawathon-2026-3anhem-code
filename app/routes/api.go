package routes

import (
	v1 "github.com/lichcse/clawathon-2026-3anhem-code/app/routes/v1"
	"github.com/lichcse/clawathon-2026-3anhem-code/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var route = gin.Default()

// SetupRouter func setup app route
func SetupRouter(db *gorm.DB) *gin.Engine {
	identityModule(db)
	postModule(db)

	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return route
}

func identityModule(db *gorm.DB) {
	iV1 := route.Group("/identity/v1")
	{
		v1.UserRoute(iV1, db)
	}
}

func postModule(db *gorm.DB) {
	pV1 := route.Group("/post/v1")
	{
		v1.PostRoute(pV1, db)
	}
}

package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"voicechat/internal/api/handler"
	"voicechat/internal/api/middleware"
	"voicechat/internal/repository"
	"voicechat/internal/service"
	"voicechat/internal/ws"
)

func NewRouter(db *sql.DB, redisClient *redis.Client, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	seatRepo := repository.NewSeatRepository(db)
	memberRepo := repository.NewRoomMemberRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtSecret)
	roomService := service.NewRoomService(roomRepo, seatRepo, memberRepo)
	seatService := service.NewSeatService(seatRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	roomHandler := handler.NewRoomHandler(roomService, seatService)
	healthHandler := handler.NewHealthHandler()

	// Health check
	router.GET("/health", healthHandler.Health)

	// Auth routes (public)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/me", middleware.AuthMiddleware(jwtSecret), authHandler.GetMe)
	}

	// Room routes
	roomGroup := router.Group("/api/v1/rooms")
	{
		roomGroup.GET("", roomHandler.ListRooms)
		roomGroup.POST("", middleware.AuthMiddleware(jwtSecret), roomHandler.CreateRoom)
		roomGroup.GET("/:room_id", roomHandler.GetRoom)
		roomGroup.DELETE("/:room_id", middleware.AuthMiddleware(jwtSecret), roomHandler.DeleteRoom)
	}

	// Seat routes
	seatGroup := router.Group("/api/v1/rooms/:room_id/seats")
	{
		seatGroup.POST("/:seat_id/occupy", middleware.AuthMiddleware(jwtSecret), roomHandler.OccupySeat)
		seatGroup.DELETE("/:seat_id", middleware.AuthMiddleware(jwtSecret), roomHandler.VacateSeat)
	}

	// WebSocket
	wsHub := ws.NewHub(redisClient)
	go wsHub.Run()

	router.GET("/ws", middleware.AuthMiddleware(jwtSecret), func(c *gin.Context) {
		userID := c.GetString("user_id")
		ws.HandleWebSocket(c.Writer, c.Request, wsHub, userID, db, memberRepo, seatRepo)
	})

	return router
}

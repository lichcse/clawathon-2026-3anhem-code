package main

import (
	"fmt"
	"log"

	"voicechat/internal/api"
	"voicechat/internal/infra"
	"voicechat/internal/infra/cache"
	"voicechat/internal/infra/db"
	"voicechat/internal/util"
)

func main() {
	// Load config
	cfg := infra.LoadConfig()
	logger := util.NewLogger()

	logger.Info("Starting Voice Chat Server...")
	logger.Infof("Environment: %s", cfg.Environment)
	logger.Infof("Server Port: %s", cfg.ServerPort)

	// Connect to database
	logger.Info("Connecting to PostgreSQL...")
	database, err := db.NewPostgresDB(cfg.GetDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Run migrations
	logger.Info("Running database migrations...")
	if err := db.RunMigrations(database); err != nil {
		log.Fatal(err)
	}

	// Connect to Redis
	logger.Info("Connecting to Redis...")
	redisClient, err := cache.NewRedisClient(cfg.RedisAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Create router
	router := api.NewRouter(database, redisClient, cfg.JWTSecret)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Infof("Server listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}

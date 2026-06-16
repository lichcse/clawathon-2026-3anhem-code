package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"voicechat/internal/api"
	"voicechat/internal/infra"
	"voicechat/internal/infra/cache"
	"voicechat/internal/infra/db"
	"voicechat/internal/infra/metrics"
	"voicechat/internal/util"
)

func main() {
	// Load config
	cfg := infra.LoadConfig()
	logger := util.NewLogger(os.Stdout)

	logger.Info("Starting Voice Chat Server...")
	logger.Infof("Environment: %s", cfg.Environment)
	logger.Infof("Server Port: %s", cfg.ServerPort)

	// Initialize Prometheus metrics
	metrics.Init()

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

	// Configure HTTP server with timeouts
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		logger.Infof("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	logger.Info("Server exited cleanly")
}

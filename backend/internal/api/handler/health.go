package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewHealthHandler(db *sql.DB, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, rdb: rdb}
}

// Live is a liveness probe: always returns 200 if the process is running.
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Ready is a readiness probe: checks DB + Redis connectivity.
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	checks := gin.H{}
	ready := true

	if err := h.db.PingContext(ctx); err != nil {
		checks["postgres"] = gin.H{"status": "unhealthy", "error": err.Error()}
		ready = false
	} else {
		checks["postgres"] = gin.H{"status": "healthy"}
	}

	if err := h.rdb.Ping(ctx).Err(); err != nil {
		checks["redis"] = gin.H{"status": "unhealthy", "error": err.Error()}
		ready = false
	} else {
		checks["redis"] = gin.H{"status": "healthy"}
	}

	status := http.StatusOK
	statusStr := "ready"
	if !ready {
		status = http.StatusServiceUnavailable
		statusStr = "not ready"
	}

	c.JSON(status, gin.H{
		"status": statusStr,
		"checks": checks,
	})
}

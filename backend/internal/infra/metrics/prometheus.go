package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	WebSocketConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_connections_active",
		Help: "Current active WebSocket connections",
	})

	WebSocketMessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_total",
			Help: "Total WebSocket messages processed",
		},
		[]string{"event"},
	)

	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	AudioFramesForwarded = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "audio_frames_forwarded_total",
		Help: "Total audio frames forwarded between clients",
	})
)

func Init() {
	prometheus.MustRegister(
		HTTPRequestDuration,
		HTTPRequestTotal,
		WebSocketConnections,
		WebSocketMessagesTotal,
		DBQueryDuration,
		AudioFramesForwarded,
	)
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		HTTPRequestDuration.WithLabelValues(c.Request.Method, c.FullPath(), status).Observe(dur)
		HTTPRequestTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
	}
}

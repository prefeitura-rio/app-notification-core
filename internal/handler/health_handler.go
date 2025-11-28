package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prefeitura-rio/app-notification-core/pkg/queue"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db       *gorm.DB
	rabbitMQ *queue.RabbitMQClient
}

func NewHealthHandler(db *gorm.DB, rabbitMQ *queue.RabbitMQClient) *HealthHandler {
	return &HealthHandler{
		db:       db,
		rabbitMQ: rabbitMQ,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status   string            `json:"status"`
	Checks   map[string]string `json:"checks,omitempty"`
	Timestamp string           `json:"timestamp"`
}

// Liveness godoc
// @Summary Liveness probe
// @Description Returns 200 if the application is running (minimal check)
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health/live [get]
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(200, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness godoc
// @Summary Readiness probe
// @Description Returns 200 if the application is ready to serve traffic (checks dependencies)
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Success 503 {object} HealthResponse
// @Router /health/ready [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	checks := make(map[string]string)
	allHealthy := true

	// Check database
	sqlDB, err := h.db.DB()
	if err != nil {
		checks["database"] = "error: " + err.Error()
		allHealthy = false
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			checks["database"] = "unreachable: " + err.Error()
			allHealthy = false
		} else {
			checks["database"] = "healthy"
		}
	}

	// Check RabbitMQ
	if h.rabbitMQ == nil {
		checks["rabbitmq"] = "not configured"
		allHealthy = false
	} else {
		_, err := h.rabbitMQ.GetQueueStats()
		if err != nil {
			checks["rabbitmq"] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			checks["rabbitmq"] = "healthy"
		}
	}

	status := "ready"
	statusCode := 200
	if !allHealthy {
		status = "not_ready"
		statusCode = 503
	}

	c.JSON(statusCode, HealthResponse{
		Status:    status,
		Checks:    checks,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Health godoc
// @Summary Basic health check
// @Description Returns 200 if the application is running (for uptime monitoring)
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(200, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

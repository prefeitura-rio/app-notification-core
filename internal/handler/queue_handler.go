package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueueMonitor interface {
	GetQueueStats() (map[string]interface{}, error)
	PurgeQueue() error
}

type QueueHandler struct {
	monitor QueueMonitor
}

func NewQueueHandler(monitor QueueMonitor) *QueueHandler {
	return &QueueHandler{monitor: monitor}
}

// GetStats godoc
// @Summary Obter estatísticas da fila
// @Description Retorna estatísticas em tempo real da fila de notificações
// @Tags queue
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /queue/stats [get]
func (h *QueueHandler) GetStats(c *gin.Context) {
	stats, err := h.monitor.GetQueueStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// PurgeQueue godoc
// @Summary Limpar fila
// @Description Remove todas as mensagens pendentes da fila
// @Tags queue
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /queue/purge [post]
func (h *QueueHandler) PurgeQueue(c *gin.Context) {
	if err := h.monitor.PurgeQueue(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Queue purged successfully"})
}

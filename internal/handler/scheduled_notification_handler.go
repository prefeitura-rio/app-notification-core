package handler

import (
	"net/http"
	"strconv"

	"github.com/fzolio/app-notification-core/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScheduledNotificationHandler struct {
	repo repository.NotificationRepository
}

func NewScheduledNotificationHandler(repo repository.NotificationRepository) *ScheduledNotificationHandler {
	return &ScheduledNotificationHandler{repo: repo}
}

// ListScheduled godoc
// @Summary Listar notificações agendadas
// @Description Retorna lista de notificações agendadas pendentes
// @Tags scheduled-notifications
// @Produce json
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {array} entity.Notification
// @Failure 500 {object} map[string]string
// @Router /scheduled-notifications [get]
func (h *ScheduledNotificationHandler) ListScheduled(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.repo.FindScheduled(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// CancelScheduled godoc
// @Summary Cancelar notificação agendada
// @Description Cancela uma notificação que estava agendada
// @Tags scheduled-notifications
// @Produce json
// @Param id path string true "ID da notificação"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /scheduled-notifications/{id}/cancel [post]
func (h *ScheduledNotificationHandler) CancelScheduled(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	if err := h.repo.CancelScheduled(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "scheduled notification cancelled"})
}

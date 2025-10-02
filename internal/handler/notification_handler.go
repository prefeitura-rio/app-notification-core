package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prefeitura-rio/app-notification-core/internal/entity"
	"github.com/prefeitura-rio/app-notification-core/internal/service"
	"github.com/prefeitura-rio/app-notification-core/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler(service service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// List godoc
// @Summary Listar notificações
// @Description Retorna lista de notificações com paginação
// @Tags notifications
// @Produce json
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {array} entity.Notification
// @Failure 500 {object} map[string]string
// @Router /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.service.ListNotifications(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetMyNotifications godoc
// @Summary Buscar minhas notificações
// @Description Retorna as notificações do usuário autenticado
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/me [get]
func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
	// Extrair informações do usuário do contexto (populado pelo middleware)
	userInfo, exists := auth.GetUserInfo(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Buscar notificações do usuário usando CPF
	notifications, err := h.service.GetNotificationsByCPF(userInfo.CPF, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar com informações do usuário
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"cpf":            userInfo.CPF,
			"email":          userInfo.Email,
			"name":           userInfo.Name,
			"email_verified": userInfo.EmailVerified,
		},
		"notifications": notifications,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(notifications),
		},
	})
}

// Create godoc
// @Summary Criar notificação
// @Description Cria uma nova notificação no sistema
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body entity.Notification true "Dados da notificação"
// @Success 201 {object} entity.Notification
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications [post]
func (h *NotificationHandler) Create(c *gin.Context) {
	var notification entity.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// Get godoc
// @Summary Buscar notificação por ID
// @Description Retorna uma notificação específica pelo ID
// @Tags notifications
// @Produce json
// @Param id path string true "ID da notificação"
// @Success 200 {object} entity.Notification
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notifications/{id} [get]
func (h *NotificationHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	notification, err := h.service.GetNotification(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// GetByCPF godoc
// @Summary Buscar notificações por CPF
// @Description Retorna notificações de um usuário específico por CPF
// @Tags notifications
// @Produce json
// @Param cpf path string true "CPF do usuário"
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {array} entity.Notification
// @Failure 500 {object} map[string]string
// @Router /notifications/cpf/{cpf} [get]
func (h *NotificationHandler) GetByCPF(c *gin.Context) {
	cpf := c.Param("cpf")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.service.GetNotificationsByCPF(cpf, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetByPhone godoc
// @Summary Buscar notificações por telefone
// @Description Retorna notificações de um usuário específico por telefone
// @Tags notifications
// @Produce json
// @Param phone path string true "Telefone do usuário"
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {array} entity.Notification
// @Failure 500 {object} map[string]string
// @Router /notifications/phone/{phone} [get]
func (h *NotificationHandler) GetByPhone(c *gin.Context) {
	phone := c.Param("phone")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.service.GetNotificationsByPhone(phone, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetByEmail godoc
// @Summary Buscar notificações por email
// @Description Retorna notificações de um usuário específico por email
// @Tags notifications
// @Produce json
// @Param email path string true "Email do usuário"
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {array} entity.Notification
// @Failure 500 {object} map[string]string
// @Router /notifications/email/{email} [get]
func (h *NotificationHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.service.GetNotificationsByEmail(email, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// Update godoc
// @Summary Atualizar notificação
// @Description Atualiza uma notificação existente
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "ID da notificação"
// @Param notification body entity.Notification true "Dados da notificação"
// @Success 200 {object} entity.Notification
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/{id} [put]
func (h *NotificationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	var notification entity.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification.ID = id
	if err := h.service.UpdateNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// Delete godoc
// @Summary Deletar notificação
// @Description Remove uma notificação do sistema
// @Tags notifications
// @Produce json
// @Param id path string true "ID da notificação"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	if err := h.service.DeleteNotification(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// MarkAsRead godoc
// @Summary Marcar notificação como lida
// @Description Marca uma notificação específica como lida
// @Tags notifications
// @Produce json
// @Param id path string true "ID da notificação"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	if err := h.service.MarkAsRead(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

type SendNotificationRequest struct {
	Title        string         `json:"title" binding:"required"`
	Message      string         `json:"message" binding:"required"`
	Type         string         `json:"type" binding:"required"`
	Data         map[string]any `json:"data,omitempty"`
	CPF          string         `json:"cpf,omitempty"`
	Phone        string         `json:"phone,omitempty"`
	Email        string         `json:"email,omitempty"`
	IsHTML       bool           `json:"is_html,omitempty"`
	IsScheduled  bool           `json:"is_scheduled,omitempty"`
	ScheduledFor *string        `json:"scheduled_for,omitempty"` // RFC3339 format
}

type BatchRecipient struct {
	CPF   string `json:"cpf,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

type SendBatchRequest struct {
	Title        string           `json:"title" binding:"required"`
	Message      string           `json:"message" binding:"required"`
	Type         string           `json:"type" binding:"required"`
	Data         map[string]any   `json:"data,omitempty"`
	IsHTML       bool             `json:"is_html,omitempty"`
	IsScheduled  bool             `json:"is_scheduled,omitempty"`
	ScheduledFor *string          `json:"scheduled_for,omitempty"` // RFC3339 format
	Recipients   []BatchRecipient `json:"recipients" binding:"required,min=1"`
}

type BatchResult struct {
	Total     int      `json:"total"`
	Succeeded int      `json:"succeeded"`
	Failed    int      `json:"failed"`
	Errors    []string `json:"errors,omitempty"`
}

// SendToUser godoc
// @Summary Enviar notificação para usuário
// @Description Envia notificação para um usuário específico via CPF, telefone ou email
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body SendNotificationRequest true "Dados da notificação"
// @Success 200 {object} entity.Notification
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/send/user [post]
func (h *NotificationHandler) SendToUser(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := &entity.Notification{
		Title:   req.Title,
		Message: req.Message,
		Type:    entity.NotificationType(req.Type),
		Data:    req.Data,
		IsHTML:  req.IsHTML,
		IsScheduled: req.IsScheduled,
	}

	// Parse scheduled_for se fornecido
	if req.IsScheduled && req.ScheduledFor != nil {
		scheduledTime, err := time.Parse(time.RFC3339, *req.ScheduledFor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_for format, use RFC3339"})
			return
		}

		// Validar se a data é futura
		if scheduledTime.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "scheduled_for must be in the future"})
			return
		}

		notification.ScheduledFor = &scheduledTime
	}

	if err := h.service.SendToUser(req.CPF, req.Phone, req.Email, notification); err != nil {
		log.Printf("Error sending notification to user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// SendToGroup godoc
// @Summary Enviar notificação para grupo
// @Description Envia notificação para todos os membros de um grupo
// @Tags notifications
// @Accept json
// @Produce json
// @Param groupId path string true "ID do grupo"
// @Param notification body SendNotificationRequest true "Dados da notificação"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/send/group/{groupId} [post]
func (h *NotificationHandler) SendToGroup(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := &entity.Notification{
		Title:       req.Title,
		Message:     req.Message,
		Type:        entity.NotificationType(req.Type),
		Data:        req.Data,
		IsHTML:      req.IsHTML,
		IsScheduled: req.IsScheduled,
	}

	// Parse scheduled_for se fornecido
	if req.IsScheduled && req.ScheduledFor != nil {
		scheduledTime, err := time.Parse(time.RFC3339, *req.ScheduledFor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_for format, use RFC3339"})
			return
		}

		// Validar se a data é futura
		if scheduledTime.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "scheduled_for must be in the future"})
			return
		}

		notification.ScheduledFor = &scheduledTime
	}

	if err := h.service.SendToGroup(groupID, notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification sent to group"})
}

// SendBroadcast godoc
// @Summary Enviar notificação broadcast
// @Description Envia notificação para todos os usuários (broadcast)
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body SendNotificationRequest true "Dados da notificação"
// @Success 200 {object} entity.Notification
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/send/broadcast [post]
func (h *NotificationHandler) SendBroadcast(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := &entity.Notification{
		Title:       req.Title,
		Message:     req.Message,
		Type:        entity.NotificationType(req.Type),
		Data:        req.Data,
		IsHTML:      req.IsHTML,
		IsScheduled: req.IsScheduled,
	}

	// Parse scheduled_for se fornecido
	if req.IsScheduled && req.ScheduledFor != nil {
		scheduledTime, err := time.Parse(time.RFC3339, *req.ScheduledFor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_for format, use RFC3339"})
			return
		}

		// Validar se a data é futura
		if scheduledTime.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "scheduled_for must be in the future"})
			return
		}

		notification.ScheduledFor = &scheduledTime
	}

	if err := h.service.SendBroadcast(notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// SendBatch godoc
// @Summary Enviar notificações em lote
// @Description Envia notificações para múltiplos destinatários em lote
// @Tags notifications
// @Accept json
// @Produce json
// @Param batch body SendBatchRequest true "Dados do envio em lote"
// @Success 200 {object} BatchResult
// @Failure 400 {object} map[string]string
// @Router /notifications/send/batch [post]
func (h *NotificationHandler) SendBatch(c *gin.Context) {
	var req SendBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding batch request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse e validar scheduled_for se fornecido
	var scheduledTime *time.Time
	if req.IsScheduled && req.ScheduledFor != nil {
		parsedTime, err := time.Parse(time.RFC3339, *req.ScheduledFor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_for format, use RFC3339"})
			return
		}

		// Validar se a data é futura
		if parsedTime.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "scheduled_for must be in the future"})
			return
		}

		scheduledTime = &parsedTime
	}

	log.Printf("Processing batch send for %d recipients", len(req.Recipients))

	result := BatchResult{
		Total:  len(req.Recipients),
		Errors: []string{},
	}

	for i, recipient := range req.Recipients {
		notification := &entity.Notification{
			Title:        req.Title,
			Message:      req.Message,
			Type:         entity.NotificationType(req.Type),
			Data:         req.Data,
			IsHTML:       req.IsHTML,
			IsScheduled:  req.IsScheduled,
			ScheduledFor: scheduledTime,
		}

		err := h.service.SendToUser(recipient.CPF, recipient.Phone, recipient.Email, notification)
		if err != nil {
			result.Failed++
			errorMsg := ""
			if recipient.Name != "" {
				errorMsg = recipient.Name + ": " + err.Error()
			} else if recipient.CPF != "" {
				errorMsg = "CPF " + recipient.CPF + ": " + err.Error()
			} else if recipient.Phone != "" {
				errorMsg = "Phone " + recipient.Phone + ": " + err.Error()
			} else if recipient.Email != "" {
				errorMsg = "Email " + recipient.Email + ": " + err.Error()
			} else {
				errorMsg = "Recipient " + string(rune(i+1)) + ": " + err.Error()
			}
			result.Errors = append(result.Errors, errorMsg)
			log.Printf("Failed to send to recipient %d: %v", i+1, err)
		} else {
			result.Succeeded++
		}
	}

	log.Printf("Batch send completed: %d succeeded, %d failed", result.Succeeded, result.Failed)

	c.JSON(http.StatusOK, result)
}

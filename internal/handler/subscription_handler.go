package handler

import (
	"net/http"

	"github.com/prefeitura-rio/app-notification-core/internal/entity"
	"github.com/prefeitura-rio/app-notification-core/internal/repository"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionHandler(repo repository.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}

type SubscribeRequest struct {
	UserCPF   string `json:"user_cpf,omitempty"`
	UserPhone string `json:"user_phone,omitempty"`
	Endpoint  string `json:"endpoint" binding:"required"`
	P256dh    string `json:"p256dh" binding:"required"`
	Auth      string `json:"auth" binding:"required"`
}

// Subscribe godoc
// @Summary Criar inscrição para push notifications
// @Description Registra uma inscrição de push notification para um usuário
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body SubscribeRequest true "Dados da inscrição"
// @Success 201 {object} entity.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription := &entity.Subscription{
		UserCPF:   req.UserCPF,
		UserPhone: req.UserPhone,
		Endpoint:  req.Endpoint,
		P256dh:    req.P256dh,
		Auth:      req.Auth,
	}

	if err := h.repo.Create(subscription); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// Unsubscribe godoc
// @Summary Cancelar inscrição de push notifications
// @Description Remove uma inscrição de push notification pelo endpoint
// @Tags subscriptions
// @Param endpoint query string true "Endpoint da inscrição"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [delete]
func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	endpoint := c.Query("endpoint")
	if endpoint == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint is required"})
		return
	}

	if err := h.repo.DeleteByEndpoint(endpoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

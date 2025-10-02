package utils

import (
	"encoding/json"
	"fmt"
	"log"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/prefeitura-rio/app-notification-core/internal/config"
	"github.com/prefeitura-rio/app-notification-core/internal/entity"
)

type WebPushClient struct {
	config *config.Config
}

func NewWebPushClient(cfg *config.Config) *WebPushClient {
	return &WebPushClient{config: cfg}
}

type PushPayload struct {
	Title   string         `json:"title"`
	Message string         `json:"message"`
	ID      string         `json:"id,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

// SendPush envia uma push notification para uma subscription
func (w *WebPushClient) SendPush(subscription *entity.Subscription, notification *entity.Notification) error {
	// Preparar payload
	payload := PushPayload{
		Title:   notification.Title,
		Message: notification.Message,
		ID:      notification.ID.String(),
	}

	if notification.Data != nil {
		payload.Data = notification.Data
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Criar subscription do webpush
	s := &webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			P256dh: subscription.P256dh,
			Auth:   subscription.Auth,
		},
	}

	// Enviar push notification
	resp, err := webpush.SendNotification(payloadBytes, s, &webpush.Options{
		Subscriber:      w.config.WebPush.VAPIDSubject,
		VAPIDPublicKey:  w.config.WebPush.VAPIDPublicKey,
		VAPIDPrivateKey: w.config.WebPush.VAPIDPrivateKey,
		TTL:             86400, // 24 horas
	})
	if err != nil {
		return fmt.Errorf("failed to send push notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("push service returned status %d", resp.StatusCode)
	}

	log.Printf("Push notification sent successfully to endpoint: %s", subscription.Endpoint)
	return nil
}

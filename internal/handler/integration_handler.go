package handler

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/fzolio/app-notification-core/internal/config"
	"github.com/gin-gonic/gin"
)

type IntegrationHandler struct {
	config *config.Config
}

func NewIntegrationHandler(cfg *config.Config) *IntegrationHandler {
	return &IntegrationHandler{config: cfg}
}

type VAPIDKeys struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Subject    string `json:"subject"`
}

type IntegrationConfig struct {
	BackendURL      string    `json:"backend_url"`
	WebSocketURL    string    `json:"websocket_url"`
	CurrentVAPID    VAPIDKeys `json:"current_vapid"`
	APIEndpoints    []string  `json:"api_endpoints"`
	SwaggerURL      string    `json:"swagger_url"`
}

// GetConfig godoc
// @Summary Obter configurações de integração
// @Description Retorna as configurações atuais para integração com frontends
// @Tags integration
// @Produce json
// @Success 200 {object} IntegrationConfig
// @Router /integration/config [get]
func (h *IntegrationHandler) GetConfig(c *gin.Context) {
	cfg := IntegrationConfig{
		BackendURL:   "http://localhost:8080/api/v1",
		WebSocketURL: "ws://localhost:8080/api/v1/ws",
		CurrentVAPID: VAPIDKeys{
			PublicKey:  h.config.WebPush.VAPIDPublicKey,
			PrivateKey: h.config.WebPush.VAPIDPrivateKey,
			Subject:    h.config.WebPush.VAPIDSubject,
		},
		APIEndpoints: []string{
			"/notifications",
			"/notifications/send/user",
			"/notifications/send/group/:id",
			"/notifications/send/batch",
			"/notifications/send/broadcast",
			"/subscriptions",
			"/groups",
			"/ws?user_id=<identifier>",
		},
		SwaggerURL: "http://localhost:8080/swagger/index.html",
	}

	c.JSON(http.StatusOK, cfg)
}

// GenerateVAPIDKeys godoc
// @Summary Gerar novas chaves VAPID
// @Description Gera um novo par de chaves VAPID para push notifications
// @Tags integration
// @Produce json
// @Success 200 {object} VAPIDKeys
// @Failure 500 {object} map[string]string
// @Router /integration/vapid/generate [post]
func (h *IntegrationHandler) GenerateVAPIDKeys(c *gin.Context) {
	// Gerar chave privada ECDSA usando curva P-256
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate private key"})
		return
	}

	// Serializar chave privada
	privateKeyBytes := privateKey.D.Bytes()
	// Padding para 32 bytes
	if len(privateKeyBytes) < 32 {
		paddedBytes := make([]byte, 32)
		copy(paddedBytes[32-len(privateKeyBytes):], privateKeyBytes)
		privateKeyBytes = paddedBytes
	}

	// Serializar chave pública (uncompressed format: 0x04 + X + Y)
	publicKeyBytes := elliptic.Marshal(elliptic.P256(), privateKey.PublicKey.X, privateKey.PublicKey.Y)

	// Codificar em base64 URL-safe (sem padding)
	publicKeyBase64 := base64.RawURLEncoding.EncodeToString(publicKeyBytes)
	privateKeyBase64 := base64.RawURLEncoding.EncodeToString(privateKeyBytes)

	keys := VAPIDKeys{
		PublicKey:  publicKeyBase64,
		PrivateKey: privateKeyBase64,
		Subject:    "mailto:your-email@example.com",
	}

	c.JSON(http.StatusOK, keys)
}

// GetEnvTemplate godoc
// @Summary Obter template de configuração
// @Description Retorna templates de .env para backend e frontend
// @Tags integration
// @Produce json
// @Success 200 {object} map[string]string
// @Router /integration/env-template [get]
func (h *IntegrationHandler) GetEnvTemplate(c *gin.Context) {
	backendEnv := `# Backend .env
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_MODE=debug

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=notification_db
DB_SSLMODE=disable

VAPID_PUBLIC_KEY=` + h.config.WebPush.VAPIDPublicKey + `
VAPID_PRIVATE_KEY=` + h.config.WebPush.VAPIDPrivateKey + `
VAPID_SUBJECT=` + h.config.WebPush.VAPIDSubject + `

DATA_RELAY_API_URL=https://data-relay.dados.rio/
DATA_RELAY_API_TOKEN=your_token_here`

	frontendEnv := `# Frontend .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_VAPID_PUBLIC_KEY=` + h.config.WebPush.VAPIDPublicKey

	c.JSON(http.StatusOK, gin.H{
		"backend":  backendEnv,
		"frontend": frontendEnv,
	})
}

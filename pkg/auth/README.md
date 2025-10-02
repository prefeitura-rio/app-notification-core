# Módulo de Autenticação JWT

Este módulo fornece utilitários para extrair informações de tokens JWT sem validar a assinatura.

**Importante**: Este módulo **NÃO valida** a assinatura do token. A validação RBAC e autenticação é feita por outra aplicação. Aqui apenas extraímos as informações do payload do JWT.

## Estrutura UserInfo

```go
type UserInfo struct {
    CPF           string   // Vem do campo preferred_username
    Email         string
    Name          string
    Phone         string
    Roles         []string
    EmailVerified bool
    Sub           string   // ID único do usuário
}
```

## Uso Básico

### 1. Parse completo do token

```go
import "github.com/prefeitura-rio/app-notification-core/pkg/auth"

func handler(c *gin.Context) {
    token := c.GetHeader("Authorization")

    userInfo, err := auth.ParseToken(token)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid token"})
        return
    }

    cpf := userInfo.CPF
    email := userInfo.Email
    name := userInfo.Name
    roles := userInfo.Roles
}
```

### 2. Extrair apenas CPF

```go
cpf, err := auth.ExtractCPF(token)
if err != nil {
    // Handle error
}
```

### 3. Extrair apenas Email

```go
email, err := auth.ExtractEmail(token)
if err != nil {
    // Handle error
}
```

## Uso com Middleware

### Proteger todas as rotas de um grupo

```go
import "github.com/prefeitura-rio/app-notification-core/pkg/auth"

func setupRoutes(router *gin.Engine) {
    v1 := router.Group("/api/v1")

    // Todas as rotas de notifications exigem autenticação
    notifications := v1.Group("/notifications")
    notifications.Use(auth.JWTMiddleware())
    {
        notifications.POST("", handler.Create)
        notifications.GET("", handler.List)
        notifications.GET("/me", handler.GetMy)
    }
}
```

### Proteger rota específica

```go
notifications.GET("/me", auth.RequireAuth(), handler.GetMyNotifications)
```

### Middleware opcional (não obrigatório)

```go
// O token será extraído se existir, mas a rota funciona sem ele
notifications.GET("/public", auth.OptionalJWTMiddleware(), handler.List)
```

### Extrair informações no handler

```go
func (h *Handler) GetMyNotifications(c *gin.Context) {
    // Extrair informações do contexto
    userInfo, exists := auth.GetUserInfo(c)
    if !exists {
        c.JSON(400, gin.H{"error": "User not authenticated"})
        return
    }

    // Usar as informações
    cpf := userInfo.CPF
    email := userInfo.Email
    name := userInfo.Name

    // Buscar notificações do usuário
    notifications, err := h.service.GetNotificationsByCPF(cpf, 20, 0)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, notifications)
}
```

## Formato do Token

O token JWT esperado segue o formato:

```
Authorization: Bearer <token>
```

Exemplo:
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Mapeamento de Campos

| Campo JWT | Campo UserInfo | Descrição |
|-----------|----------------|-----------|
| `preferred_username` | `CPF` | CPF do usuário |
| `email` | `Email` | Email do usuário |
| `name` | `Name` | Nome completo |
| `phone_number` | `Phone` | Telefone |
| `email_verified` | `EmailVerified` | Email verificado |
| `sub` | `Sub` | ID único do usuário |
| `realm_access.roles` | `Roles` | Lista de roles |

## Exemplo Completo

```go
package handler

import (
    "github.com/prefeitura-rio/app-notification-core/pkg/auth"
    "github.com/gin-gonic/gin"
)

type NotificationHandler struct {
    service NotificationService
}

// GetMyNotifications retorna as notificações do usuário autenticado
func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
    // Extrair informações do usuário do contexto (populado pelo middleware)
    userInfo, exists := auth.GetUserInfo(c)
    if !exists {
        c.JSON(401, gin.H{"error": "Unauthorized"})
        return
    }

    // Usar CPF para buscar notificações
    notifications, err := h.service.GetNotificationsByCPF(userInfo.CPF, 20, 0)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "user": gin.H{
            "cpf":   userInfo.CPF,
            "email": userInfo.Email,
            "name":  userInfo.Name,
        },
        "notifications": notifications,
    })
}

// Configurar rotas
func SetupRoutes(router *gin.Engine, handler *NotificationHandler) {
    v1 := router.Group("/api/v1")

    notifications := v1.Group("/notifications")
    {
        // Rota pública (sem autenticação)
        notifications.GET("/public", handler.ListPublic)

        // Rotas protegidas (requerem autenticação)
        notifications.GET("/me", auth.RequireAuth(), handler.GetMyNotifications)
        notifications.POST("/me/read/:id", auth.RequireAuth(), handler.MarkAsRead)
    }
}
```

## Tratamento de Erros

O middleware retorna os seguintes erros:

- `401 Unauthorized` - Token ausente ou inválido
- O token é considerado inválido se:
  - Não estiver no formato JWT (3 partes separadas por '.')
  - O payload não puder ser decodificado de Base64
  - O JSON do payload for inválido

## Notas Importantes

1. **Sem Validação de Assinatura**: Este módulo não valida a assinatura do token. Assume-se que a validação já foi feita por outra aplicação.

2. **Extração de Dados**: O módulo apenas extrai e parseia o payload do JWT para uso na aplicação.

3. **CPF**: O CPF vem do campo `preferred_username` do token IDRio.

4. **Middleware Opcional**: Use `OptionalJWTMiddleware()` quando a autenticação for opcional mas você quiser as informações do usuário se disponíveis.

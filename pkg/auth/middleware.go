package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserInfoKey = "user_info"
)

// JWTMiddleware extrai informações do token JWT e adiciona ao contexto
// Este middleware NÃO valida a assinatura do token (isso é feito por outra aplicação)
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai o token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		// Parse do token
		userInfo, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Adiciona as informações do usuário ao contexto
		c.Set(UserInfoKey, userInfo)

		c.Next()
	}
}

// OptionalJWTMiddleware é similar ao JWTMiddleware mas não retorna erro se o token não existir
// Útil para rotas que podem funcionar com ou sem autenticação
func OptionalJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			token = strings.TrimSpace(token)

			userInfo, err := ParseToken(token)
			if err == nil {
				c.Set(UserInfoKey, userInfo)
			}
		}

		c.Next()
	}
}

// GetUserInfo extrai as informações do usuário do contexto
func GetUserInfo(c *gin.Context) (*UserInfo, bool) {
	userInfo, exists := c.Get(UserInfoKey)
	if !exists {
		return nil, false
	}

	user, ok := userInfo.(*UserInfo)
	return user, ok
}

// RequireAuth garante que o usuário está autenticado
func RequireAuth() gin.HandlerFunc {
	return JWTMiddleware()
}

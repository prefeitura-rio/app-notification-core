package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// UserInfo contém as informações extraídas do token JWT
type UserInfo struct {
	CPF           string   `json:"cpf"`
	Email         string   `json:"email"`
	Name          string   `json:"name"`
	Phone         string   `json:"phone"`
	Roles         []string `json:"roles"`
	EmailVerified bool     `json:"email_verified"`
	Sub           string   `json:"sub"`
}

// JWTClaims representa os claims do token JWT
type JWTClaims struct {
	PreferredUsername string   `json:"preferred_username"`
	Email             string   `json:"email"`
	Name              string   `json:"name"`
	EmailVerified     bool     `json:"email_verified"`
	Phone             string   `json:"phone_number"`
	Sub               string   `json:"sub"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

// ParseToken faz o parse de um token JWT e extrai as informações do usuário
// Nota: Esta função NÃO valida a assinatura do token, apenas extrai os dados
func ParseToken(tokenString string) (*UserInfo, error) {
	// Remove "Bearer " se presente
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	// Split do token em partes
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	// Decode do payload (segunda parte)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("failed to decode token payload")
	}

	// Parse dos claims
	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("failed to parse token claims")
	}

	// Extrai informações do usuário
	userInfo := &UserInfo{
		CPF:           claims.PreferredUsername, // CPF vem no preferred_username
		Email:         claims.Email,
		Name:          claims.Name,
		Phone:         claims.Phone,
		EmailVerified: claims.EmailVerified,
		Sub:           claims.Sub,
		Roles:         claims.RealmAccess.Roles,
	}

	return userInfo, nil
}

// ExtractCPF extrai apenas o CPF do token
func ExtractCPF(tokenString string) (string, error) {
	userInfo, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return userInfo.CPF, nil
}

// ExtractEmail extrai apenas o email do token
func ExtractEmail(tokenString string) (string, error) {
	userInfo, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return userInfo.Email, nil
}

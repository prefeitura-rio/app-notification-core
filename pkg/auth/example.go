package auth

// Exemplo de uso do middleware e funções de autenticação
//
// 1. Aplicar middleware em rotas específicas:
//
//   notifications := v1.Group("/notifications")
//   notifications.Use(auth.JWTMiddleware()) // Todas as rotas exigem autenticação
//   {
//       notifications.POST("", notificationHandler.Create)
//       notifications.GET("", notificationHandler.List)
//   }
//
// 2. Aplicar middleware em rota específica:
//
//   notifications.GET("/me", auth.RequireAuth(), notificationHandler.GetMyNotifications)
//
// 3. Usar middleware opcional (não obrigatório):
//
//   notifications.GET("/public", auth.OptionalJWTMiddleware(), notificationHandler.List)
//
// 4. Extrair informações do usuário no handler:
//
//   func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
//       userInfo, exists := auth.GetUserInfo(c)
//       if !exists {
//           c.JSON(400, gin.H{"error": "User info not found"})
//           return
//       }
//
//       cpf := userInfo.CPF
//       email := userInfo.Email
//       name := userInfo.Name
//
//       // Buscar notificações do usuário
//       notifications, err := h.service.GetNotificationsByCPF(cpf, 20, 0)
//       // ...
//   }
//
// 5. Extrair CPF diretamente do token (sem middleware):
//
//   token := c.GetHeader("Authorization")
//   cpf, err := auth.ExtractCPF(token)
//   if err != nil {
//       c.JSON(400, gin.H{"error": "Invalid token"})
//       return
//   }

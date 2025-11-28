package main

import (
	"log"
	"time"

	_ "github.com/prefeitura-rio/app-notification-core/docs"
	"github.com/prefeitura-rio/app-notification-core/internal/config"
	"github.com/prefeitura-rio/app-notification-core/internal/handler"
	"github.com/prefeitura-rio/app-notification-core/internal/repository"
	"github.com/prefeitura-rio/app-notification-core/internal/scheduler"
	"github.com/prefeitura-rio/app-notification-core/internal/service"
	"github.com/prefeitura-rio/app-notification-core/internal/websocket"
	"github.com/prefeitura-rio/app-notification-core/pkg/auth"
	"github.com/prefeitura-rio/app-notification-core/pkg/queue"
	"github.com/prefeitura-rio/app-notification-core/pkg/utils"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Notification Service API
// @version 1.0
// @description API para gerenciamento de notificações com suporte a WebSocket e Push Notifications
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Token JWT no formato: Bearer {token}

func main() {
	// Configurar timezone para horário de Brasília
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Printf("Warning: Failed to load America/Sao_Paulo timezone: %v. Using system default.", err)
	} else {
		time.Local = loc
		log.Printf("Timezone set to: %s", loc.String())
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	groupRepo := repository.NewGroupRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	hub := websocket.NewHub()
	go hub.Run()

	mailman := utils.NewMailmanClient(cfg.DataRelay.URL, cfg.DataRelay.Token)
	webPush := utils.NewWebPushClient(cfg)

	// Conectar ao RabbitMQ
	rabbitMQ, err := queue.NewRabbitMQClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	groupService := service.NewGroupService(groupRepo)
	notificationService := service.NewNotificationService(notificationRepo, groupRepo, subscriptionRepo, hub, mailman, webPush, rabbitMQ)

	// Iniciar scheduler de notificações agendadas
	notificationScheduler := scheduler.NewNotificationScheduler(notificationRepo, notificationService)
	notificationScheduler.Start()
	defer notificationScheduler.Stop()

	// Iniciar workers para processar a fila
	workers := cfg.RabbitMQ.Workers
	if workers == 0 {
		workers = 3 // Default
	}
	log.Printf("Starting %d workers to process notifications...", workers)
	for i := 0; i < workers; i++ {
		workerID := i + 1
		go func(id int) {
			log.Printf("Worker %d started", id)
			rabbitMQ.ConsumeNotifications(func(msg *queue.NotificationMessage) error {
				return notificationService.ProcessNotification(msg.Notification)
			})
		}(workerID)
	}

	groupHandler := handler.NewGroupHandler(groupService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	scheduledNotificationHandler := handler.NewScheduledNotificationHandler(notificationRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionRepo)
	wsHandler := handler.NewWebSocketHandler(hub)
	integrationHandler := handler.NewIntegrationHandler(cfg)
	queueHandler := handler.NewQueueHandler(rabbitMQ)
	healthHandler := handler.NewHealthHandler(db, rabbitMQ)

	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	router.Use(corsMiddleware())

	v1 := router.Group("/api/v1")
	{
		groups := v1.Group("/groups")
		{
			groups.POST("", groupHandler.Create)
			groups.GET("", groupHandler.List)
			groups.GET("/:id", groupHandler.Get)
			groups.PUT("/:id", groupHandler.Update)
			groups.DELETE("/:id", groupHandler.Delete)
			groups.POST("/:id/members", groupHandler.AddMember)
			groups.GET("/:id/members", groupHandler.GetMembers)
			groups.GET("/:id/members/:memberId", groupHandler.GetMember)
			groups.PUT("/:id/members/:memberId", groupHandler.UpdateMember)
			groups.DELETE("/:id/members/:memberId", groupHandler.RemoveMember)
		}

		notifications := v1.Group("/notifications")
		{
			notifications.POST("", notificationHandler.Create)
			notifications.GET("", notificationHandler.List)

			// Rotas autenticadas (requerem JWT)
			notifications.GET("/me", auth.RequireAuth(), notificationHandler.GetMyNotifications)

			notifications.GET("/:id", notificationHandler.Get)
			notifications.PUT("/:id", notificationHandler.Update)
			notifications.DELETE("/:id", notificationHandler.Delete)
			notifications.POST("/:id/read", notificationHandler.MarkAsRead)

			notifications.GET("/cpf/:cpf", notificationHandler.GetByCPF)
			notifications.GET("/phone/:phone", notificationHandler.GetByPhone)
			notifications.GET("/email/:email", notificationHandler.GetByEmail)

			notifications.POST("/send/user", notificationHandler.SendToUser)
			notifications.POST("/send/group/:groupId", notificationHandler.SendToGroup)
			notifications.POST("/send/broadcast", notificationHandler.SendBroadcast)
			notifications.POST("/send/batch", notificationHandler.SendBatch)
		}

		scheduledNotifications := v1.Group("/scheduled-notifications")
		{
			scheduledNotifications.GET("", scheduledNotificationHandler.ListScheduled)
			scheduledNotifications.POST("/:id/cancel", scheduledNotificationHandler.CancelScheduled)
		}

		v1.GET("/ws", wsHandler.ServeWS)

		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.Subscribe)
			subscriptions.DELETE("", subscriptionHandler.Unsubscribe)
		}

		integration := v1.Group("/integration")
		{
			integration.GET("/config", integrationHandler.GetConfig)
			integration.POST("/vapid/generate", integrationHandler.GenerateVAPIDKeys)
			integration.GET("/env-template", integrationHandler.GetEnvTemplate)
		}

		queue := v1.Group("/queue")
		{
			queue.GET("/stats", queueHandler.GetStats)
			queue.POST("/purge", queueHandler.PurgeQueue)
		}
	}

	// Health check endpoints
	router.GET("/health", healthHandler.Health)
	router.GET("/health/live", healthHandler.Liveness)
	router.GET("/health/ready", healthHandler.Readiness)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

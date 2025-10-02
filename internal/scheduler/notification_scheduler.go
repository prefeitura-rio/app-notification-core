package scheduler

import (
	"log"
	"time"

	"github.com/prefeitura-rio/app-notification-core/internal/entity"
	"github.com/prefeitura-rio/app-notification-core/internal/repository"
	"github.com/prefeitura-rio/app-notification-core/internal/service"
)

type NotificationScheduler struct {
	notificationRepo repository.NotificationRepository
	notificationService service.NotificationService
	ticker *time.Ticker
	stopChan chan bool
}

func NewNotificationScheduler(
	repo repository.NotificationRepository,
	service service.NotificationService,
) *NotificationScheduler {
	return &NotificationScheduler{
		notificationRepo: repo,
		notificationService: service,
		stopChan: make(chan bool),
	}
}

// Start inicia o scheduler que verifica notificações agendadas a cada minuto
func (s *NotificationScheduler) Start() {
	log.Println("📅 Notification Scheduler started")

	// Processar imediatamente ao iniciar
	s.processScheduledNotifications()

	// Processar a cada 1 minuto
	s.ticker = time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.processScheduledNotifications()
			case <-s.stopChan:
				log.Println("📅 Notification Scheduler stopped")
				return
			}
		}
	}()
}

// Stop para o scheduler
func (s *NotificationScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopChan)
}

// processScheduledNotifications busca e processa notificações agendadas que estão prontas para envio
func (s *NotificationScheduler) processScheduledNotifications() {
	now := time.Now()

	// Buscar notificações agendadas cuja hora já passou
	notifications, err := s.notificationRepo.FindScheduledReady(now)
	if err != nil {
		log.Printf("❌ Error fetching scheduled notifications: %v", err)
		return
	}

	if len(notifications) == 0 {
		return
	}

	log.Printf("📅 Found %d scheduled notification(s) ready to send", len(notifications))

	for _, notification := range notifications {
		go s.sendScheduledNotification(&notification)
	}
}

// sendScheduledNotification envia uma notificação agendada
func (s *NotificationScheduler) sendScheduledNotification(notification *entity.Notification) {
	log.Printf("📤 Sending scheduled notification: %s (ID: %s)", notification.Title, notification.ID)

	// Atualizar status para pending (vai ser processado pelo worker)
	notification.Status = entity.StatusPending
	notification.IsScheduled = false
	if err := s.notificationRepo.Update(notification); err != nil {
		log.Printf("❌ Failed to update scheduled notification %s: %v", notification.ID, err)
		return
	}

	// Enviar notificação (será processada pela fila)
	if err := s.notificationService.SendNotification(notification); err != nil {
		log.Printf("❌ Failed to send scheduled notification %s: %v", notification.ID, err)
		s.notificationRepo.UpdateStatus(notification.ID, entity.StatusFailed)
		return
	}

	log.Printf("✅ Scheduled notification sent: %s", notification.ID)
}

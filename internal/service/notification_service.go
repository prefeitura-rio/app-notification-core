package service

import (
	"errors"
	"log"

	"github.com/prefeitura-rio/app-notification-core/internal/entity"
	"github.com/prefeitura-rio/app-notification-core/internal/repository"
	"github.com/prefeitura-rio/app-notification-core/internal/websocket"
	"github.com/prefeitura-rio/app-notification-core/pkg/utils"
	"github.com/google/uuid"
)

type NotificationService interface {
	CreateNotification(notification *entity.Notification) error
	GetNotification(id uuid.UUID) (*entity.Notification, error)
	ListNotifications(limit, offset int) ([]entity.Notification, error)
	GetNotificationsByCPF(cpf string, limit, offset int) ([]entity.Notification, error)
	GetNotificationsByPhone(phone string, limit, offset int) ([]entity.Notification, error)
	GetNotificationsByEmail(email string, limit, offset int) ([]entity.Notification, error)
	UpdateNotification(notification *entity.Notification) error
	DeleteNotification(id uuid.UUID) error
	MarkAsRead(id uuid.UUID) error
	SendNotification(notification *entity.Notification) error
	ProcessNotification(notification *entity.Notification) error
	SendToUser(cpf, phone, email string, notification *entity.Notification) error
	SendToGroup(groupID uuid.UUID, notification *entity.Notification) error
	SendBroadcast(notification *entity.Notification) error
}

type notificationService struct {
	notificationRepo   repository.NotificationRepository
	groupRepo          repository.GroupRepository
	subscriptionRepo   repository.SubscriptionRepository
	hub                *websocket.Hub
	mailman            *utils.MailmanClient
	webPush            *utils.WebPushClient
	queue              QueuePublisher
}

// QueuePublisher interface para publicar mensagens na fila
type QueuePublisher interface {
	PublishNotification(notification *entity.Notification) error
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	groupRepo repository.GroupRepository,
	subscriptionRepo repository.SubscriptionRepository,
	hub *websocket.Hub,
	mailman *utils.MailmanClient,
	webPush *utils.WebPushClient,
	queue QueuePublisher,
) NotificationService {
	return &notificationService{
		notificationRepo:   notificationRepo,
		groupRepo:          groupRepo,
		subscriptionRepo:   subscriptionRepo,
		hub:                hub,
		mailman:            mailman,
		webPush:            webPush,
		queue:              queue,
	}
}

func (s *notificationService) CreateNotification(notification *entity.Notification) error {
	if notification.Title == "" || notification.Message == "" {
		return errors.New("title and message are required")
	}
	return s.notificationRepo.Create(notification)
}

func (s *notificationService) GetNotification(id uuid.UUID) (*entity.Notification, error) {
	return s.notificationRepo.FindByID(id)
}

func (s *notificationService) ListNotifications(limit, offset int) ([]entity.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.notificationRepo.FindAll(limit, offset)
}

func (s *notificationService) GetNotificationsByCPF(cpf string, limit, offset int) ([]entity.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.notificationRepo.FindByCPF(cpf, limit, offset)
}

func (s *notificationService) GetNotificationsByPhone(phone string, limit, offset int) ([]entity.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.notificationRepo.FindByPhone(phone, limit, offset)
}

func (s *notificationService) GetNotificationsByEmail(email string, limit, offset int) ([]entity.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.notificationRepo.FindByEmail(email, limit, offset)
}

func (s *notificationService) UpdateNotification(notification *entity.Notification) error {
	if notification.Title == "" || notification.Message == "" {
		return errors.New("title and message are required")
	}
	return s.notificationRepo.Update(notification)
}

func (s *notificationService) DeleteNotification(id uuid.UUID) error {
	return s.notificationRepo.Delete(id)
}

func (s *notificationService) MarkAsRead(id uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(id)
}

func (s *notificationService) SendNotification(notification *entity.Notification) error {
	log.Printf("SendNotification: Creating notification with type=%s", notification.Type)

	// Verificar se é uma notificação agendada
	if notification.IsScheduled && notification.ScheduledFor != nil {
		notification.Status = entity.StatusScheduled
		if err := s.CreateNotification(notification); err != nil {
			log.Printf("SendNotification: Failed to create scheduled notification: %v", err)
			return err
		}
		log.Printf("SendNotification: Scheduled notification %s created for %s",
			notification.ID, notification.ScheduledFor.Format("2006-01-02 15:04:05"))
		return nil
	}

	// Notificação imediata
	if err := s.CreateNotification(notification); err != nil {
		log.Printf("SendNotification: Failed to create notification: %v", err)
		return err
	}

	// Publicar na fila RabbitMQ para processamento assíncrono
	if err := s.queue.PublishNotification(notification); err != nil {
		log.Printf("SendNotification: Failed to publish to queue: %v", err)
		s.notificationRepo.UpdateStatus(notification.ID, entity.StatusFailed)
		return err
	}

	log.Printf("SendNotification: Notification %s published to queue", notification.ID)
	return nil
}

// ProcessNotification processa a notificação (chamado pelos workers)
func (s *notificationService) ProcessNotification(notification *entity.Notification) error {
	log.Printf("ProcessNotification: Processing notification %s with type=%s", notification.ID, notification.Type)

	shouldSendInApp := notification.Type == entity.TypeInApp ||
		notification.Type == entity.TypeBoth ||
		notification.Type == entity.TypeAll

	shouldSendPush := notification.Type == entity.TypePush ||
		notification.Type == entity.TypeBoth ||
		notification.Type == entity.TypeAll

	shouldSendEmail := notification.Type == entity.TypeEmail ||
		notification.Type == entity.TypeAll

	log.Printf("ProcessNotification: shouldSendInApp=%v, shouldSendPush=%v, shouldSendEmail=%v", shouldSendInApp, shouldSendPush, shouldSendEmail)

	if shouldSendInApp {
		log.Printf("ProcessNotification: Broadcasting notification via WebSocket")
		s.hub.BroadcastNotification(notification)
	}

	if shouldSendPush {
		log.Printf("ProcessNotification: Sending push notifications")
		s.sendPushNotifications(notification)
	}

	if shouldSendEmail && notification.UserEmail != nil && *notification.UserEmail != "" {
		log.Printf("ProcessNotification: Sending email to %s", *notification.UserEmail)
		mailReq := &utils.MailmanRequest{
			ToAddresses: []string{*notification.UserEmail},
			Subject:     notification.Title,
			Body:        notification.Message,
			IsHTMLBody:  notification.IsHTML,
		}

		if err := s.mailman.SendEmail(mailReq); err != nil {
			log.Printf("ProcessNotification: Failed to send email: %v", err)
			s.notificationRepo.UpdateStatus(notification.ID, entity.StatusFailed)
			return err
		}
		log.Printf("ProcessNotification: Email sent successfully")
	}

	if err := s.notificationRepo.UpdateStatus(notification.ID, entity.StatusSent); err != nil {
		log.Printf("ProcessNotification: Failed to update status: %v", err)
		return err
	}

	log.Printf("ProcessNotification: Notification %s processed successfully", notification.ID)
	return nil
}

// sendPushNotifications envia push notifications para as subscriptions do usuário
func (s *notificationService) sendPushNotifications(notification *entity.Notification) {
	var subscriptions []entity.Subscription
	var err error

	// Buscar subscriptions baseado no identificador disponível
	if notification.UserCPF != nil && *notification.UserCPF != "" {
		subscriptions, err = s.subscriptionRepo.FindByCPF(*notification.UserCPF)
		if err != nil {
			log.Printf("Failed to find subscriptions by CPF: %v", err)
			return
		}
	} else if notification.UserPhone != nil && *notification.UserPhone != "" {
		subscriptions, err = s.subscriptionRepo.FindByPhone(*notification.UserPhone)
		if err != nil {
			log.Printf("Failed to find subscriptions by phone: %v", err)
			return
		}
	}

	if len(subscriptions) == 0 {
		log.Printf("No subscriptions found for notification %s", notification.ID)
		return
	}

	log.Printf("Found %d subscription(s), sending push notifications...", len(subscriptions))

	// Enviar push notification para cada subscription
	for _, sub := range subscriptions {
		if err := s.webPush.SendPush(&sub, notification); err != nil {
			log.Printf("Failed to send push to subscription %s: %v", sub.ID, err)
			// Continuar enviando para outras subscriptions mesmo se uma falhar
			continue
		}
		log.Printf("Push sent successfully to subscription %s", sub.ID)
	}
}

func (s *notificationService) SendToUser(cpf, phone, email string, notification *entity.Notification) error {
	if cpf != "" {
		notification.UserCPF = &cpf
	}
	if phone != "" {
		notification.UserPhone = &phone
	}
	if email != "" {
		notification.UserEmail = &email
	}

	return s.SendNotification(notification)
}

func (s *notificationService) SendToGroup(groupID uuid.UUID, notification *entity.Notification) error {
	members, err := s.groupRepo.FindMembers(groupID)
	if err != nil {
		return err
	}

	notification.GroupID = &groupID

	for _, member := range members {
		individualNotif := *notification
		individualNotif.ID = uuid.Nil
		if member.CPF != "" {
			individualNotif.UserCPF = &member.CPF
		}
		if member.Phone != "" {
			individualNotif.UserPhone = &member.Phone
		}
		if member.Email != "" {
			individualNotif.UserEmail = &member.Email
		}

		if err := s.SendNotification(&individualNotif); err != nil {
			continue
		}
	}

	return nil
}

func (s *notificationService) SendBroadcast(notification *entity.Notification) error {
	notification.Broadcast = true
	return s.SendNotification(notification)
}

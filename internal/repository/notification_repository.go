package repository

import (
	"time"

	"github.com/fzolio/app-notification-core/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *entity.Notification) error
	FindByID(id uuid.UUID) (*entity.Notification, error)
	FindAll(limit, offset int) ([]entity.Notification, error)
	FindByCPF(cpf string, limit, offset int) ([]entity.Notification, error)
	FindByPhone(phone string, limit, offset int) ([]entity.Notification, error)
	FindByEmail(email string, limit, offset int) ([]entity.Notification, error)
	FindByGroupID(groupID uuid.UUID, limit, offset int) ([]entity.Notification, error)
	Update(notification *entity.Notification) error
	Delete(id uuid.UUID) error
	MarkAsRead(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status entity.NotificationStatus) error
	FindScheduledReady(before time.Time) ([]entity.Notification, error)
	FindScheduled(limit, offset int) ([]entity.Notification, error)
	CancelScheduled(id uuid.UUID) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *entity.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByID(id uuid.UUID) (*entity.Notification, error) {
	var notification entity.Notification
	err := r.db.First(&notification, "id = ?", id).Error
	return &notification, err
}

func (r *notificationRepository) FindAll(limit, offset int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindByCPF(cpf string, limit, offset int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("user_cpf = ? OR broadcast = ?", cpf, true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindByPhone(phone string, limit, offset int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("user_phone = ? OR broadcast = ?", phone, true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindByEmail(email string, limit, offset int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("user_email = ? OR broadcast = ?", email, true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindByGroupID(groupID uuid.UUID, limit, offset int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("group_id = ?", groupID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) Update(notification *entity.Notification) error {
	return r.db.Save(notification).Error
}

func (r *notificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Notification{}, "id = ?", id).Error
}

func (r *notificationRepository) MarkAsRead(id uuid.UUID) error {
	now := gorm.Expr("NOW()")
	return r.db.Model(&entity.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  entity.StatusRead,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) UpdateStatus(id uuid.UUID, status entity.NotificationStatus) error {
	return r.db.Model(&entity.Notification{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// FindScheduledReady busca notificações agendadas prontas para envio
func (r *notificationRepository) FindScheduledReady(before time.Time) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("is_scheduled = ? AND scheduled_for <= ? AND status = ?",
		true, before, entity.StatusScheduled).
		Order("scheduled_for ASC").
		Find(&notifications).Error
	return notifications, err
}

// FindScheduled lista todas as notificações agendadas
func (r *notificationRepository) FindScheduled(limit, offset int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.db.Where("is_scheduled = ? AND status = ?", true, entity.StatusScheduled).
		Order("scheduled_for ASC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// CancelScheduled cancela uma notificação agendada
func (r *notificationRepository) CancelScheduled(id uuid.UUID) error {
	return r.db.Model(&entity.Notification{}).
		Where("id = ? AND is_scheduled = ? AND status = ?", id, true, entity.StatusScheduled).
		Updates(map[string]interface{}{
			"status":      entity.StatusCancelled,
			"is_scheduled": false,
		}).Error
}

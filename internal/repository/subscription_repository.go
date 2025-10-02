package repository

import (
	"github.com/prefeitura-rio/app-notification-core/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *entity.Subscription) error
	FindByEndpoint(endpoint string) (*entity.Subscription, error)
	FindByCPF(cpf string) ([]entity.Subscription, error)
	FindByPhone(phone string) ([]entity.Subscription, error)
	Delete(id uuid.UUID) error
	DeleteByEndpoint(endpoint string) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(subscription *entity.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *subscriptionRepository) FindByEndpoint(endpoint string) (*entity.Subscription, error) {
	var subscription entity.Subscription
	err := r.db.First(&subscription, "endpoint = ?", endpoint).Error
	return &subscription, err
}

func (r *subscriptionRepository) FindByCPF(cpf string) ([]entity.Subscription, error) {
	var subscriptions []entity.Subscription
	err := r.db.Where("user_cpf = ?", cpf).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) FindByPhone(phone string) ([]entity.Subscription, error) {
	var subscriptions []entity.Subscription
	err := r.db.Where("user_phone = ?", phone).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Subscription{}, "id = ?", id).Error
}

func (r *subscriptionRepository) DeleteByEndpoint(endpoint string) error {
	return r.db.Delete(&entity.Subscription{}, "endpoint = ?", endpoint).Error
}

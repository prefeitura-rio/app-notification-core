package entity

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string
type NotificationStatus string

const (
	TypeInApp NotificationType = "in-app"
	TypePush  NotificationType = "push"
	TypeEmail NotificationType = "email"
	TypeBoth  NotificationType = "both"
	TypeAll   NotificationType = "all"

	StatusPending   NotificationStatus = "pending"
	StatusSent      NotificationStatus = "sent"
	StatusDelivered NotificationStatus = "delivered"
	StatusRead      NotificationStatus = "read"
	StatusFailed    NotificationStatus = "failed"
)

type Notification struct {
	ID        uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey"`
	Title     string             `json:"title" gorm:"not null"`
	Message   string             `json:"message" gorm:"not null"`
	Type      NotificationType   `json:"type" gorm:"not null"`
	Status    NotificationStatus `json:"status" gorm:"default:'pending'"`
	Data      map[string]any     `json:"data,omitempty" gorm:"type:jsonb"`
	UserCPF   *string            `json:"user_cpf,omitempty" gorm:"index"`
	UserPhone *string            `json:"user_phone,omitempty" gorm:"index"`
	UserEmail *string            `json:"user_email,omitempty" gorm:"index"`
	GroupID   *uuid.UUID         `json:"group_id,omitempty" gorm:"type:uuid;index"`
	Broadcast bool               `json:"broadcast" gorm:"default:false"`
	IsHTML    bool               `json:"is_html" gorm:"default:false"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	ReadAt    *time.Time         `json:"read_at,omitempty"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

package entity

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserCPF   string    `json:"user_cpf" gorm:"index"`
	UserPhone string    `json:"user_phone" gorm:"index"`
	Endpoint  string    `json:"endpoint" gorm:"not null;uniqueIndex"`
	P256dh    string    `json:"p256dh" gorm:"not null"`
	Auth      string    `json:"auth" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

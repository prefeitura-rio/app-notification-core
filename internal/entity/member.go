package entity

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	GroupID   uuid.UUID `json:"group_id" gorm:"type:uuid;not null"`
	CPF       string    `json:"cpf" gorm:"index"`
	Phone     string    `json:"phone" gorm:"index"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

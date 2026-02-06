// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/freepass-2026/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	OrderID   uuid.UUID      `json:"order_id" gorm:"type:char(36);"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:char(36);"`
	Content   string         `json:"content" gorm:"type:varchar(1024)"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (f *Feedback) ParseToDTOResponseCreateFeedback() dto.ResponseCreateFeedback {
	return dto.ResponseCreateFeedback{
		ID:        f.ID,
		OrderID:   f.OrderID,
		UserID:    f.UserID,
		Content:   f.Content,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func (f *Feedback) ParseToDTOResponseGetFeedback() dto.ResponseGetFeedback {
	return dto.ResponseGetFeedback{
		ID:        f.ID,
		OrderID:   f.OrderID,
		UserID:    f.UserID,
		Content:   f.Content,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

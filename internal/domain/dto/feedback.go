// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateFeedback struct {
	ID        uuid.UUID      `json:"id" validate:"required,uuid_rfc4122"`
	OrderID   uuid.UUID      `json:"order_id" validate:"required,uuid_rfc4122"`
	UserID    uuid.UUID      `json:"user_id" validate:"required,uuid_rfc4122"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ResponseCreateFeedback struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SoftDeleteFeedback struct {
	ID uuid.UUID `json:"id" validate:"required,uuid_rfc4122"`
}

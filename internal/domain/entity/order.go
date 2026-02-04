// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	CanteenID uuid.UUID      `json:"order_id" gorm:"type:char(36);"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:char(36);"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

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
	MenuID    uuid.UUID      `json:"menu_id" gorm:"type:char(36);"`
	Quantity  uint32         `json:"quantity" gorm:"type:uint32"`
	Status    string         `json:"status" gorm:"type:varchar(128)"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

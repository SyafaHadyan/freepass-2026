// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Food struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	CanteenID uuid.UUID      `json:"canteen_id" gorm:"type:char(36)"`
	Name      string         `json:"name" gorm:"type:varchar(128)"`
	Price     uint32         `json:"price" gorm:"type:uint32"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

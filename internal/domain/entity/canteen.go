// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Canteen struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:char(36)"`
	Name      string         `json:"name" gorm:"tpe:varchar(128)"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

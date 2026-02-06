// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/freepass-2026/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Canteen struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:char(36)"`
	Name      string         `json:"name" gorm:"type:varchar(128)"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (c *Canteen) ParseToDTOResponseCreateCanteen() dto.ResponseCreateCanteen {
	return dto.ResponseCreateCanteen{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (c *Canteen) ParseToDTOResponseGetCanteenList() dto.ResponseGetCanteenList {
	return dto.ResponseGetCanteenList{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (c *Canteen) ParseToDTOResponseGetCanteenInfo() dto.ResponseGetCanteenInfo {
	return dto.ResponseGetCanteenInfo{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// Package entity defines database table and its relations
package entity

import (
	"time"

	"github.com/SyafaHadyan/freepass-2026/internal/domain/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	CanteenID uuid.UUID      `json:"canteen_id" gorm:"type:char(36)"`
	Name      string         `json:"name" gorm:"type:varchar(128)"`
	Price     uint32         `json:"price" gorm:"type:uint32"`
	Stock     uint32         `json:"stock" gorm:"type:uint32"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Menu) ParseToDTOResponseCreateMenu() dto.ResponseCreateMenu {
	return dto.ResponseCreateMenu{
		ID:        m.ID,
		CanteenID: m.CanteenID,
		Name:      m.Name,
		Price:     m.Price,
		Stock:     m.Stock,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m *Menu) ParseToDTOResponseUpdateMenu() dto.ResponseUpdateMenu {
	return dto.ResponseUpdateMenu{
		ID:        m.ID,
		CanteenID: m.CanteenID,
		Name:      m.Name,
		Price:     m.Price,
		Stock:     m.Stock,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m *Menu) ParseToDTOResponseGetMenuInfo() dto.ResponseGetMenuInfo {
	return dto.ResponseGetMenuInfo{
		ID:        m.ID,
		CanteenID: m.CanteenID,
		Name:      m.Name,
		Price:     m.Price,
		Stock:     m.Stock,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

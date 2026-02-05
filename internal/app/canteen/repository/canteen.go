// Package repository handles the CRUD operations
package repository

import (
	"github.com/SyafaHadyan/freepass-2026/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CanteenDBItf interface {
	CreateCanteen(canteen *entity.Canteen) error
	CreateMenu(menu *entity.Menu) error
	UpdateMenu(menu *entity.Menu, userID uuid.UUID) error
	GetCanteenInfo(canteen *entity.Canteen) error
	GetCanteenList(canteen *[]entity.Canteen) error
	GetMenuInfo(menu *entity.Menu) error
	SoftDeleteMenu(menu *entity.Menu, userID uuid.UUID) error
}

type CanteenDB struct {
	db *gorm.DB
}

func NewCanteenDB(db *gorm.DB) CanteenDBItf {
	return &CanteenDB{
		db: db,
	}
}

func (r *CanteenDB) CreateCanteen(canteen *entity.Canteen) error {
	return r.db.Debug().
		Create(canteen).
		Error
}

// TODO: validate canteen ownership

func (r *CanteenDB) CreateMenu(menu *entity.Menu) error {
	return r.db.Debug().
		Create(menu).
		Error
}

func (r *CanteenDB) UpdateMenu(menu *entity.Menu, userID uuid.UUID) error {
	sub := r.db.Debug().
		Model(&entity.Canteen{}).
		Select("id").
		Where("user_id = ?", userID)

	res := r.db.Debug().
		Where("id = ?", menu.ID).
		Where("canteen_id IN (?)", sub).
		Updates(menu)

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return res.Error
}

func (r *CanteenDB) GetCanteenList(canteen *[]entity.Canteen) error {
	return r.db.Debug().
		Model(&canteen).
		Select("id, name").
		Error
}

func (r *CanteenDB) GetCanteenInfo(canteen *entity.Canteen) error {
	return r.db.Debug().
		Select("id, user_id, name, created_at, updated_at").
		First(canteen).
		Error
}

func (r *CanteenDB) GetMenuInfo(menu *entity.Menu) error {
	return r.db.Debug().
		Select("id, canteen_id, name, price, created_at, updated_at").
		First(&menu).
		Error
}

func (r *CanteenDB) SoftDeleteMenu(menu *entity.Menu, userID uuid.UUID) error {
	sub := r.db.Debug().
		Model(&entity.Canteen{}).
		Select("id").
		Where("user_id = ?", userID)

	res := r.db.Debug().
		Where("id = ?", menu.ID).
		Where("canteen_id IN (?)", sub).
		Delete(menu)

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return res.Error
}

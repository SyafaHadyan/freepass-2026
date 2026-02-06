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
	CreateOrder(menu *entity.Menu, order *entity.Order) error
	CreatePayment(payment *entity.Payment) error
	CreateFeedback(feedback *entity.Feedback) error
	UpdateMenu(menu *entity.Menu, userID uuid.UUID) error
	UpdateOrder(order *entity.Order, userID uuid.UUID) error
	GetCanteenInfo(canteen *entity.Canteen) error
	GetCanteenList(canteen *[]entity.Canteen) error
	GetMenuInfo(menu *entity.Menu) error
	GetOrderInfo(order *entity.Order) error
	GetOrderList(order *[]entity.Order, userID uuid.UUID) error
	SoftDeleteMenu(menu *entity.Menu, userID uuid.UUID) error
	SoftDeleteFeedback(feedback *entity.Feedback, userID uuid.UUID) error
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

func (r *CanteenDB) CreateOrder(menu *entity.Menu, order *entity.Order) error {
	sub := r.db.Debug().
		Model(&menu).
		Select("id, canteen_id").
		Where("stock >= ?", order.Quantity)
	if sub == nil {
		return gorm.ErrInvalidValue
	}

	order.CanteenID = menu.CanteenID

	r.db.Debug().
		Create(order)

	return r.db.Debug().
		Model(&menu).
		Update("stock = ?", menu.Stock-order.Quantity).
		Error
}

func (r *CanteenDB) CreatePayment(payment *entity.Payment) error {
	return r.db.Debug().
		Create(payment).
		Error
}

func (r *CanteenDB) CreateFeedback(feedback *entity.Feedback) error {
	res := r.db.Debug().
		Model(&entity.Order{}).
		Where("id = ?", feedback.OrderID).
		Where("user_id = ?", feedback.UserID).
		Where("status = ?", "COMPLETED").
		Update("status", "FEEDBACKSENT")

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.db.Debug().
		Create(feedback).
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

func (r *CanteenDB) UpdateOrder(order *entity.Order, userID uuid.UUID) error {
	sub := r.db.Debug().
		Model(&entity.Canteen{}).
		Select("id").
		Where("user_id = ?", userID).
		Where("status = ?", "PAID").
		Or("status = ?", "COOKING")

	res := r.db.Debug().
		Where("id = ?", order.ID).
		Where("canteen_id IN (?)", sub).
		Update("status = ?", order.Status)

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

func (r *CanteenDB) GetOrderInfo(order *entity.Order) error {
	return r.db.Debug().
		Select("id, canteen_id, user_id, menu_id, quantity, status, created_at, updated_at").
		Where("user_id = ?", order.UserID).
		Error
}

func (r *CanteenDB) GetOrderList(order *[]entity.Order, userID uuid.UUID) error {
	sub := r.db.Debug().
		Model(&entity.Canteen{}).
		Select("id").
		Where("user_id = ?", userID)

	res := r.db.Debug().
		Model(&entity.Order{}).
		Select("id, canteen_id, user_id, menu_id, quantity, status, created_at, updated_at").
		Where("canteen_id IN (?)", sub).
		Find(order)

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return res.Error
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

func (r *CanteenDB) SoftDeleteFeedback(feedback *entity.Feedback, userID uuid.UUID) error {
	canteenSub := r.db.Debug().
		Model(&entity.Canteen{}).
		Select("id").
		Where("user_id = ?", userID)

	orderSub := r.db.Debug().
		Model(&entity.Order{}).
		Select("id").
		Where("status = ?", "FEEDBACKSENT").
		Where("canteen_id IN (?)", canteenSub)

	res := r.db.Debug().
		Where("id = ?", feedback.ID).
		Where("order_id IN (?)", orderSub).
		Delete(feedback)

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return res.Error
}

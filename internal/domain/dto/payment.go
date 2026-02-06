// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreatePayment struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	OrderID     uuid.UUID      `json:"order_id" gorm:"type:char(36);"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:char(36);"`
	Price       uint32         `json:"price" gorm:"type:uint32"`
	RedirectURL string         `json:"redirect_url" gorm:"type:varchar(256)"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type CreateMidtransOrder struct {
	TransactionDetails TransactionDetails
	CustomerDetail     CustomerDetail
}

type TransactionDetails struct {
	OrderID     string `json:"order_id" validate:"required,number,min=1"`
	GrossAmount uint32 `json:"gross_amount" validate:"required,number,min=1"`
}

type CustomerDetail struct {
	FirstName string `json:"first_name" validate:"omitempty,min=1"`
	LastName  string `json:"last_name" validate:"omitempty,min=1"`
	Email     string `json:"email" validate:"omitempty,email"`
}

type ResponseMidtransOrder struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

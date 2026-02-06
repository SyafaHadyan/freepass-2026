// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMenu struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id" validate:"required,uuid_rfc4122"`
	Name      string    `json:"name" validate:"required,min=3,max=64"`
	Price     uint32    `json:"price" validate:"required,number,min=1"`
	Stock     uint32    `json:"stock" validate:"required,number,min=1"`
}

type ResponseCreateMenu struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id"`
	Name      string    `json:"name"`
	Price     uint32    `json:"price"`
	Stock     uint32    `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateMenu struct {
	ID     uuid.UUID `json:"id" validate:"required,required,uuid_rfc4122"`
	UserID uuid.UUID `json:"user_id" validate:"required,required,uuid_rfc4122"`
	Name   string    `json:"name" validate:"omitempty,min=3,max=64"`
	Price  uint32    `json:"price" validate:"omitempty,number,min=1"`
	Stock  uint32    `json:"stock" validate:"omitempty,number,min=1"`
}

type ResponseUpdateMenu struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id"`
	Name      string    `json:"name"`
	Price     uint32    `json:"price"`
	Stock     uint32    `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetMenuInfo struct {
	ID        uuid.UUID `json:"id"`
	CanteenID uuid.UUID `json:"canteen_id"`
	Name      string    `json:"name"`
	Price     uint32    `json:"price"`
	Stock     uint32    `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SoftDeleteMenu struct {
	ID uuid.UUID `json:"id" validate:"required,required,uuid_rfc4122"`
}

// Package dto defines standarized struct to be used as data exchange
package dto

import "github.com/google/uuid"

type CreateOrder struct {
	ID        uuid.UUID `json:"id" validate:"required,uuid_rfc4122"`
	CanteenID uuid.UUID `json:"canteen_id" validate:"required,uuid_rfc4122"`
	UserID    uuid.UUID `json:"user_id" validate:"required,uuid_rfc4122"`
	MenuID    uuid.UUID `json:"menu_id" validate:"required,uuid_rfc4122"`
	Quantity  uint32    `json:"quantity" validate:"required,number,min=1"`
	Status    string    `json:"status" validate:"required,oneof=UNPAID PAID COOKING COMPLETED"`
}

// Package dto defines standarized struct to be used as data exchange
package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCanteen struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name" validate:"required,min=3,max=64"`
}

type ResponseCreateCanteen struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseGetCanteenList struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ResponseGetCanteenInfo struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

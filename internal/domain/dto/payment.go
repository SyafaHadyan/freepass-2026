// Package dto defines standarized struct to be used as data exchange
package dto

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

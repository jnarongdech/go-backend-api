package dto

import "github.com/google/uuid"

type CreateCustomerRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Phone       string `json:"phone"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	City        string `json:"City"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
}

type UpdateCustomerRequest struct {
	ID          uuid.UUID `json:"-" params:"id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Email       string    `json:"email" validate:"required"`
	Phone       string    `json:"phone"`
	CompanyName string    `json:"company_name"`
	Address     string    `json:"address"`
	City        string    `json:"City"`
	PostalCode  string    `json:"postal_code"`
	Country     string    `json:"country"`
}

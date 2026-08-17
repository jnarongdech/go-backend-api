package dto

import "encoding/json"

type CreateProductRequest struct {
	Name                string          `json:"name" validate:"required"`
	Description         *string         `json:"description"`
	Category            *string         `json:"category"`
	BasePrice           string          `json:"base_price"`
	IsCustomizable      *bool           `json:"is_customizable"`
	CustomizationFields json.RawMessage `json:"customization_fields" swaggertype:"object"`
}

type ProductResponse struct {
	ID                  string          `json:"id"`
	Name                string          `json:"name"`
	Description         *string         `json:"description"`
	Category            *string         `json:"category"`
	BasePrice           string          `json:"base_price"`
	IsCustomizable      *bool           `json:"is_customizable"`
	CustomizationFields json.RawMessage `json:"customization_fields" swaggertype:"object"`
}

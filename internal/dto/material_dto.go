package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMaterialRequest struct {
	Name           string  `json:"name" validate:"required"`
	ThicknessMM    float64 `json:"thickness_mm"`
	Grade          string  `json:"grade"`
	Description    string  `json:"description"`
	CostPerKg      float64 `json:"cost_per_kg" validate:"required,gt=0"`
	StockQtyKg     float64 `json:"stock_qty_kg"`
	ReorderLevelKg float64 `json:"reorder_level_kg"`
}

type MaterialResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	ThicknessMM    float64   `json:"thickness_mm,omitempty"`
	Grade          string    `json:"grade,omitempty"`
	Description    string    `json:"description,omitempty"`
	CostPerKg      float64   `json:"cost_per_kg"`
	StockQtyKg     float64   `json:"stock_qty_kg"`
	ReorderLevelKg float64   `json:"reorder_level_kg,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

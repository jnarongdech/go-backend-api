package dto

import (
	"github.com/google/uuid"
)

type CreateOrderItemsRequest struct {
	OrderID      uuid.UUID `json:"order_id" validate:"required"`
	ProductID    uuid.UUID `json:"product_id" validate:"required"`
	Quantity     int32     `json:"quantity" validate:"required,min=1"`
	PricePerUnit float64   `json:"price_per_unit" validate:"required,gt=0"`
}

type OrderItemsResponse struct {
	ID           uuid.UUID `json:"item_id"`
	ProductID    uuid.UUID `json:"product_id"`
	Quantity     int32     `json:"quantity"`
	PricePerUnit float64   `json:"price_per_unit"`
}

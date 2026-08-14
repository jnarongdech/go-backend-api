package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	CustomerID             uuid.UUID                 `json:"customer_id" validate:"required"`
	OrderType              string                    `json:"order_type" validate:"required,oneof=pre_order stock"`
	Status                 string                    `json:"status" validate:"omitempty,oneof=draft pending_payment paid in_production production_complete ready_to_ship shipped delivered cancelled"`
	TotalPrice             *float64                  `json:"total_price"`
	Notes                  *string                   `json:"notes"`
	OrderDate              *time.Time                `json:"order_date"`
	ExpectedCompletionDate *time.Time                `json:"expected_completion_date"`
	ActualCompletionDate   *time.Time                `json:"actual_completion_date"`
	Items                  []CreateOrderItemsRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateOrderRequest struct {
	ID                     uuid.UUID                 `json:"-" params:"id" validate:"required"`
	CustomerID             uuid.UUID                 `json:"customer_id" validate:"required"`
	OrderType              string                    `json:"order_type" validate:"required,oneof=pre_order stock"`
	Status                 string                    `json:"status" validate:"omitempty,oneof=draft pending_payment paid in_production production_complete ready_to_ship shipped delivered cancelled"`
	TotalPrice             *float64                  `json:"total_price"`
	Notes                  *string                   `json:"notes"`
	OrderDate              *time.Time                `json:"order_date"`
	ExpectedCompletionDate *time.Time                `json:"expected_completion_date"`
	ActualCompletionDate   *time.Time                `json:"actual_completion_date"`
	Items                  []CreateOrderItemsRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderResponse struct {
	ID                     uuid.UUID            `json:"id"`
	OrderNumber            string               `json:"order_number"`
	CustomerID             uuid.UUID            `json:"customer_id" `
	OrderType              interface{}          `json:"order_type"`
	Status                 interface{}          `json:"status"`
	TotalPrice             *float64             `json:"total_price"`
	Notes                  *string              `json:"notes"`
	OrderDate              *time.Time           `json:"order_date"`
	ExpectedCompletionDate *time.Time           `json:"expected_completion_date"`
	ActualCompletionDate   *time.Time           `json:"actual_completion_date"`
	Items                  []OrderItemsResponse `json:"items"`
}

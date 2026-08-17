package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jnarongdech/go-backend-api/internal/dto"
	"github.com/jnarongdech/go-backend-api/internal/service"
	constants "github.com/jnarongdech/go-backend-api/pkg/consts"
	"github.com/jnarongdech/go-backend-api/pkg/errs"
	"github.com/jnarongdech/go-backend-api/pkg/response"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// GetOrdersHandler godoc
// @Summary      Get all data...
// @Description  Retrieves a list of items from the system.
// @Tags         Orders
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /api/v1/orders [get]
func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	data, err := h.orderService.GetOrders(c.Context())
	if err != nil {
		return err
	}

	orderList := make([]dto.OrderResponse, 0, len(data))
	for _, dataItem := range data {
		mappedItem := dto.OrderResponse{
			ID:                     dataItem.ID,
			OrderNumber:            dataItem.OrderNumber,
			CustomerID:             dataItem.CustomerID,
			OrderType:              dataItem.OrderType,
			Status:                 dataItem.Status,
			TotalPrice:             dataItem.TotalPrice,
			Notes:                  dataItem.Notes,
			OrderDate:              dataItem.OrderDate,
			ExpectedCompletionDate: dataItem.ExpectedCompletionDate,
			ActualCompletionDate:   dataItem.ActualCompletionDate,
		}
		orderList = append(orderList, mappedItem)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Get orders successfully!", orderList)
}

// GetProductByIDHandler godoc
// @Summary      Get data by ID...
// @Description  Retrieves a single item by its ID.
// @Tags         Orders
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200 {object} map[string]interface{}
// @Router       /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	orderID := c.Params("id")
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return errs.NewBadRequest(constants.ErrInvalidIDFormat, err)
	}

	result, err := h.orderService.GetOrderByID(c.Context(), orderUUID)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Success!", result)
}

// CreateOrderWithItemsHandler godoc
// @Summary      Create new data...
// @Description  Takes a JSON payload and creates an item in the system.
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateOrderRequest true "Data to create"
// @Success      201 {object} map[string]interface{}
// @Router       /api/v1/orders [post]
func (h *OrderHandler) CreateOrderWithItems(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequest(constants.ErrInvalidJSONFormat, err)
	}

	err := h.orderService.CreateOrderWithItems(c.Context(), req)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Created Order Successfully!", nil)
}

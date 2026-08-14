package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jnarongdech/go-backend-api/internal/dto"
	"github.com/jnarongdech/go-backend-api/internal/repository"
	constants "github.com/jnarongdech/go-backend-api/pkg/consts"
	"github.com/jnarongdech/go-backend-api/pkg/errs"
	"github.com/jnarongdech/go-backend-api/pkg/response"
)

type ICustomerService interface {
	GetCustomerByID(ctx context.Context, id uuid.UUID) (repository.Customer, error)
	CreateCustomer(ctx context.Context, req dto.CreateCustomerRequest) (repository.Customer, error)
	UpdateCustomer(ctx context.Context, req dto.UpdateCustomerRequest) error
}

type CustomerHandler struct {
	customerService ICustomerService
}

// Constructor Function
func NewCustomerHandler(customerService ICustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

// GetCustomerByIDHandler godoc
// @Summary      Get data by ID...
// @Description  Retrieves a single item by its ID.
// @Tags         Customers
// @Produce      json
// @Param        id path string true "Customer ID"
// @Success      200 {object} map[string]interface{}
// @Router       /api/v1/customers/{id} [get]
func (h *CustomerHandler) GetCustomerByID(c *fiber.Ctx) error {
	// GET ID from URL PATH /api/v1/customers/:id
	customerID := c.Params("id")
	if customerID == "" {
		return response.ErrorResponse(c, errors.New("Customer ID is required"))
	}
	customerUUID, err := uuid.Parse(customerID)
	customer, err := h.customerService.GetCustomerByID(c.Context(), customerUUID)
	if err != nil {
		return response.ErrorResponse(c, errors.New("Customer not found."))
	}

	return response.SuccessResponse(c, 200, "Success!", customer)
}

// CreateCustomerHandler godoc
// @Summary      Create new data...
// @Description  Takes a JSON payload and creates an item in the system.
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateCustomerRequest true "Data to create"
// @Success      201 {object} map[string]interface{}
// @Router       /api/v1/customers [post]
func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req dto.CreateCustomerRequest

	if err := c.BodyParser(&req); err != nil {
		return response.ErrorResponse(c, errors.New(constants.ErrInvalidJSONFormat))
	}

	if req.Name == "" {
		return response.ErrorResponse(c, errors.New(constants.ErrMissingFieldName))
	}
	if req.Email == "" {
		return response.ErrorResponse(c, errors.New(constants.ErrMissingFieldEmail))
	}

	customer, errCreate := h.customerService.CreateCustomer(c.Context(), req)
	if errCreate != nil {
		return response.ErrorResponse(c, errCreate)
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Customer created successfully!", customer)
}

// UpdateCustomerHandler godoc
// @Summary      Update data...
// @Description  Updates an existing item by ID.
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Param        request body dto.UpdateCustomerRequest true "Data to update"
// @Success      200 {object} map[string]interface{}
// @Router       /api/v1/customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	idParam := c.Params("id")
	customerUUID, errUUID := uuid.Parse(idParam)
	if errUUID != nil {
		fmt.Println("[Error] %w", errUUID)
		return response.ErrorResponse(c, errs.NewBadRequest(constants.ErrInvalidIDFormat, errUUID))
	}

	var req dto.UpdateCustomerRequest
	req.ID = customerUUID
	if errParams := c.BodyParser(&req); errParams != nil {
		return response.ErrorResponse(c, errors.New(constants.ErrInvalidJSONFormat))
	}
	if req.Name == "" {
		return response.ErrorResponse(c, errors.New(constants.ErrMissingFieldName))
	}
	if req.Email == "" {
		return response.ErrorResponse(c, errors.New(constants.ErrMissingFieldEmail))
	}

	errUpdate := h.customerService.UpdateCustomer(c.Context(), req)
	if errUpdate != nil {
		return response.ErrorResponse(c, errUpdate)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Customer updated successfully!", nil)
}

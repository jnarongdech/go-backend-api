package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jnarongdech/go-backend-api/internal/dto"
	"github.com/jnarongdech/go-backend-api/internal/service"
	constants "github.com/jnarongdech/go-backend-api/pkg/consts"
	"github.com/jnarongdech/go-backend-api/pkg/response"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// GetProductsHandler godoc
// @Summary      Get all products...
// @Description  Retrieves a list of all products ordered from newest to oldest.
// @Tags         Products
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /api/v1/products [get]
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetProducts(c.Context())
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, "Success!", products)
}

// CreateProductHandler godoc
// @Summary Create new data...
// @Description  Takes a JSON payload and creates a product in the system.
// @Tags Products
// @Accept json
// @Param request body dto.CreateProductRequest true "Data to create"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if errBody := c.BodyParser(&req); errBody != nil {
		return response.ErrorResponse(c, errors.New(constants.ErrCreateInternalServer))
	}

	result, err := h.productService.CreateProduct(c.Context(), req)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Created Successfully!", result)
}

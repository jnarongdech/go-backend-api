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

type MaterialHandler struct {
	service service.MaterialService
}

func NewMaterialHandler(service service.MaterialService) *MaterialHandler {
	return &MaterialHandler{service: service}
}

func (h *MaterialHandler) CreateMaterial(c *fiber.Ctx) error {
	var req dto.CreateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequest(constants.ErrInvalidJSONFormat, err)
	}

	result, err := h.service.CreateMaterial(c.Context(), req)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Created Successfully!", result)
}

func (h *MaterialHandler) GetMaterials(c *fiber.Ctx) error {
	result, err := h.service.GetMaterails(c.Context())
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, "Data retrieved successfully.", result)
}

func (h *MaterialHandler) GetMaterialByID(c *fiber.Ctx) error {
	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return errs.NewBadRequest(constants.ErrInvalidIDFormat, err)
	}

	result, err := h.service.GetMaterialByID(c.Context(), uuid)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Data retrieved successfully.", result)
}

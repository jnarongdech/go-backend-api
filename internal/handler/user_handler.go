package handler

import (
	"net/mail"

	"github.com/gofiber/fiber/v2"
	"github.com/jnarongdech/go-backend-api/internal/service"
	constants "github.com/jnarongdech/go-backend-api/pkg/consts"
	"github.com/jnarongdech/go-backend-api/pkg/errs"
	"github.com/jnarongdech/go-backend-api/pkg/response"
)

// สร้าง Struct เพื่อรับ JSON จากฝั่งผู้ใช้งาน
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
}

// สำหรับ แก้ไขผู้ใช้งาน (สังเกตว่าเปลี่ยนเป็น Pointer *string เพื่อให้ส่งค่าว่าง หรือไม่ส่งมาก็ได้)
type UpdateUserRequest struct {
	ID       string `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
}

type SoftDeleteRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfileHandler godoc
// @Summary      Get data by ID...
// @Description  Retrieves a single item by its ID.
// @Tags         Users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]interface{}
// @Router       /api/v1/products/{id} [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// 1. ดึง userID จาก Context (ถ้าผ่าน Auth Middleware มา) หรือจาก Query/Param
	// userID := c.Locals("userID").(string)

	// ดึงจาก URL Param เช่น /api/v1/users/me?id=xxx
	userID := c.Query("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	// 2. เรียกใช้ Service (ซึ่งจะไปเรียก repo -> sqlc ต่ออีกที)
	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// 3. ส่ง Response กลับเป็น JSON
	return c.JSON(fiber.Map{
		"message": "Success",
		"data":    user,
	})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// สร้างตัวแปรมารอรับข้อมูล
	var req CreateUserRequest

	// ใช้ BodyParser เพื่อแกะ JSON แปลงเข้า Struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": 400,
			"error":  "Invalid email format.",
		})
	}

	// validator email
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": 400,
			"error":  "Invalid email format.",
		})
	}

	// ส่งข้อมูลที่แกะแล้ว ไปให้ Service จัดการต่อ
	user, err := h.userService.CreateUser(c.Context(), req.Email, req.FullName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": 400,
			"error":  err.Error(),
		})
	}

	// ถ้าสำเร็จ ส่งข้อมูลผู้ใช้งานใหม่กลับไปให้ Frontend
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  200,
		"success": true,
		"message": "Created Successfully!",
		"data":    user,
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequest(constants.ErrInvalidJSONFormat, err)
	}

	if err := validate.Struct(req); err != nil {
		return errs.NewBadRequest(constants.ErrInvalidFormatOrIncompleteData, err)
	}

	user, err := h.userService.UpdateUser(c.Context(), req.ID, req.Email, req.FullName)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Updated!", user)
}

func (h *UserHandler) SoftDeleteUser(c *fiber.Ctx) error {
	var req SoftDeleteRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequest(constants.ErrInvalidJSONFormat, err)
	}

	if req.ID != "" {
		return errs.NewBadRequest("Invalid ID format (UUID required).", nil)
	}

	errSoftDel := h.userService.SoftDeleteUser(c.Context(), req.ID)
	if errSoftDel != nil {
		return errSoftDel
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Disabled!", nil)
}

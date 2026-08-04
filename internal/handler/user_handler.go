package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jnarongdech/go-backend-api/internal/service"
)

// สร้าง Struct เพื่อรับ JSON จากฝั่งผู้ใช้งาน
type CreateUserRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

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
	// 1. สร้างตัวแปรมารอรับข้อมูล
	var req CreateUserRequest

	// 2. ใช้ BodyParser เพื่อแกะ JSON แปลงเข้า Struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบข้อมูลไม่ถูกต้อง",
		})
	}

	// 3. ส่งข้อมูลที่แกะแล้ว ไปให้ Service จัดการต่อ
	user, err := h.userService.CreateUser(c.Context(), req.ID, req.Email, req.FullName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 4. ถ้าสำเร็จ ส่งข้อมูลผู้ใช้งานใหม่กลับไปให้ Frontend
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "สร้างผู้ใช้งานสำเร็จ!",
		"data":    user,
	})
}

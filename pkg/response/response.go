package response

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jnarongdech/go-backend-api/pkg/errs" // เปลี่ยน path ให้ตรงกับโปรเจกต์
)

// Success Response (ตัวเดิมของคุณ ไม่ต้องแก้)
func SuccessResponse(c *fiber.Ctx, status int, msg string, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"success": true,
		"message": msg,
		"data":    data,
	})
}

// เปลี่ยนชื่อจาก ErrorResponse เป็น ErrorHandler
// สังเกตว่ารับค่า (c *fiber.Ctx, err error) ตามมาตรฐาน Fiber เป๊ะ
func ErrorHandler(c *fiber.Ctx, err error) error {
	// 1. ตั้งค่า Default
	statusCode := fiber.StatusInternalServerError
	errMsg := "An error occurred, please try again."

	// 2. เช็คว่าเป็น AppError ของเราไหม
	var appErr errs.AppError
	if errors.As(err, &appErr) {
		statusCode = appErr.Code
		errMsg = appErr.Message
		if appErr.Err != nil {
			log.Printf("[AppError] %v: %v", appErr.Message, appErr.Err)
		}
	} else {
		// ดักจับ Error พื้นฐานของ Fiber เองด้วย
		// (เช่น กรณีพิมพ์ URL ผิด Fiber จะพ่น 404 ออกมา เราก็จับมาแปลงเป็น JSON ซะเลย)
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			statusCode = fiberErr.Code
			errMsg = fiberErr.Message
		} else {
			log.Printf("[UNEXPECTED ERROR] %v", err)
		}
	}

	// 3. พ่น JSON สวยๆ กลับไป
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"success": false,
		"error":   errMsg,
	})
}

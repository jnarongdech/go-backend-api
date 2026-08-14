package response

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jnarongdech/go-backend-api/pkg/errs" // อย่าลืมเปลี่ยน path ให้ตรงกับโปรเจกต์คุณนะครับ
)

func SuccessResponse(c *fiber.Ctx, status int, msg string, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"success": true,
		"message": msg,
		"data":    data,
	})
}

// ปรับให้รับค่า err เป็น type error ตรงๆ
func ErrorResponse(c *fiber.Ctx, err error) error {
	// ตั้งค่า Default เผื่อไว้ก่อน (เผื่อเป็น error ทั่วไปที่ไม่ใช่ AppError)
	statusCode := fiber.StatusInternalServerError
	errMsg := "An error occurred, please try again."

	// ใช้ errors.As เพื่อเช็คว่า err ที่ส่งมา คือ AppError หรือเปล่า?
	var appErr errs.AppError
	if errors.As(err, &appErr) {
		// ถ้าใช่ ให้เอา Code กับ Message ที่เราตั้งไว้มาใช้
		statusCode = appErr.Code
		errMsg = appErr.Message

		// (Optional) พิมพ์ Log ดิบให้โปรแกรมเมอร์ดูหลังบ้าน
		if appErr.Err != nil {
			log.Printf("[AppError] %v: %v", appErr.Message, appErr.Err)
		}
	} else {
		// ถ้าไม่ใช่ AppError แปลว่าอาจจะหลุดมาจากไลบรารีอื่น ให้พ่น Log ตัวแดงไว้
		log.Printf("[UNEXPECTED ERROR] %v", err)
	}

	// ส่ง JSON กลับไปให้ Front-end
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode, // ใส่เลข Status ให้ตรงกับ HTTP Status
		"success": false,
		"error":   errMsg, // ข้อความสวยๆ ที่ User อ่านรู้เรื่อง
	})
}

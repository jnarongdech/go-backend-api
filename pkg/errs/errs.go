package errs

// AppError คือ Struct กลางที่เราจะใช้ทั้งโปรเจกต์
type AppError struct {
	Code    int    // HTTP Status Code (เช่น 400, 404, 500)
	Message string // ข้อความสวยๆ ที่จะส่งให้ Front-end ไปโชว์
	Err     error  // Error ดิบๆ เอาไว้ Log ลง Console (จะไม่ส่งไปหา Front-end)
}

// สำคัญมาก: ต้องสร้างฟังก์ชัน Error() เพื่อให้ AppError สวมรอยเป็น error ของ Go ได้
func (e AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// ฟังก์ชันช่วยเหลือ (Helper) เพื่อให้เรียกใช้ได้สั้นๆ
func NewBadRequest(message string, err error) error {
	return AppError{
		Code:    400,
		Message: message,
		Err:     err,
	}
}

func NewInternal(message string, err error) error {
	return AppError{
		Code:    500,
		Message: message,
		Err:     err,
	}
}

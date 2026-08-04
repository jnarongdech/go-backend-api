package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jnarongdech/go-backend-api/internal/repository"
)

// 1. สร้าง Struct เก็บ Repository ไว้ใช้งาน
type UserService struct {
	userRepo *repository.UserRepository
}

// 2. ฟังก์ชัน Constructor (เรียกใช้ใน main.go)
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// 3. เริ่มเขียน Business Logic ต่างๆ ให้ Handler เรียกใช้งาน
// ฟังก์ชันดึงข้อมูล User
func (s *UserService) GetUserByID(ctx context.Context, id string) (repository.User, error) {
	// 1. (Business Logic): ตรวจสอบความถูกต้องเบื้องต้น
	if id == "" {
		// ถ้าไม่มี ID ส่งมา ให้ตีกลับทันที
		return repository.User{}, errors.New("user id cannot be empty")
	}

	// เรียกใช้งาน Repository เพื่อไปดึงข้อมูลจาก Database
	// ฟังก์ชัน FindUserByID นี้คือตัวที่เราเขียนจัดการแปลง String เป็น UUID ไว้ใน repo
	user, err := s.userRepo.FindUserByID(ctx, id)

	if err != nil {
		// 1. เก็บ Log ข้อความ Error ของจริง (Raw Error) เอาไว้ให้เราดูเอง
		// เวลาเกิดปัญหา เราจะได้มาเปิดดูในหน้าจอ Terminal หรือระบบ Log ได้ว่าพังเพราะอะไร
		log.Printf("[ERROR] GetUserByID failed for id %s: %v", id, err)

		// 2. แปลงข้อความ Error ให้เป็นมิตรและปลอดภัยก่อนส่งให้ Frontend
		// เช็คว่า Error ที่เกิดขึ้น คือการ "หาข้อมูลไม่เจอ" (No Rows) ใช่หรือไม่
		if errors.Is(err, sql.ErrNoRows) {
			// ถ้าใช่ ให้ส่งข้อความที่เข้าใจง่ายกลับไป
			return repository.User{}, errors.New("ไม่พบข้อมูลผู้ใช้งานในระบบ")
		}

		// ถ้าเป็น Error แบบอื่นๆ (เช่น Database ล่ม, เน็ตหลุด) เราจะไม่บอกสาเหตุจริงๆ ให้ผู้ใช้รู้
		// แต่จะส่งข้อความกลางๆ กลับไปแทนเพื่อความปลอดภัย
		return repository.User{}, errors.New("เกิดข้อผิดพลาดขัดข้องภายในระบบ โปรดลองใหม่อีกครั้ง")
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, id string, email string, fullName string) (repository.User, error) {
	// กฎข้อที่ 1: ต้องใส่อีเมลเสมอ
	if email == "" {
		return repository.User{}, errors.New("กรุณาระบุอีเมล")
	}

	// จัดการกับตัวแปรที่อนุญาตให้เป็น Null ได้ (full_name)
	var fullNamePtr *string
	if fullName != "" {
		fullNamePtr = &fullName // ถ้ามีคนกรอกชื่อมา ให้ชี้ Pointer ไปที่ค่านั้น
	}

	// ส่งข้อมูลไปให้ Repository จัดการ
	user, err := s.userRepo.CreateUser(ctx, id, email, fullNamePtr)
	if err != nil {
		log.Printf("[ERROR] CreateUser failed: %v", err)
		return repository.User{}, errors.New("ไม่สามารถสร้างผู้ใช้งานได้ในขณะนี้")
	}

	return user, nil
}

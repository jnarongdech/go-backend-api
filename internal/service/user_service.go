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
		return repository.User{}, errors.New("เกิดข้อผิดพลาดขัดข้อง โปรดลองใหม่อีกครั้ง")
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, email string, fullName string) (repository.User, error) {
	if email == "" || fullName == "" {
		return repository.User{}, errors.New("email and fullname cannot an empty")
	}
	// ส่งข้อมูลไปให้ Repository จัดการ
	user, err := s.userRepo.CreateUser(ctx, email, fullName)
	if err != nil {
		log.Printf("[ERROR] CreateUser failed: %v", err)
		return repository.User{}, errors.New("ไม่สามารถสร้างข้อมูลผู้ใช้งานได้")
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, email string, fullName string) (repository.User, error) {
	if id == "" || email == "" || fullName == "" {
		return repository.User{}, errors.New("id, email and fullname cannot be empty")
	}

	user, err := s.userRepo.UpdateUser(ctx, id, email, fullName)
	if err != nil {
		log.Printf("[Error] UpdateUser failed: %v", err)
		return repository.User{}, errors.New("Unable to edit user data")
	}

	return user, nil
}

func (s *UserService) SoftDeleteUser(ctx context.Context, id string) error {
	err := s.userRepo.SoftDeleteUser(ctx, id)
	if err != nil {
		log.Printf("[Error] SoftDeleteUser failed: %v", err)
		return errors.New("Unable disable user data")
	}

	return nil
}

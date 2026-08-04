package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// สมมติว่าไฟล์ที่ sqlc สร้างให้ (models.go, query.sql.go) อยู่ใน package นี้แล้ว

type UserRepository struct {
	q *Queries
}

func NewUserRepository(q *Queries) *UserRepository {
	return &UserRepository{q: q}
}

// รับ id มาเป็น string เหมือนเดิม
func (r *UserRepository) FindUserByID(ctx context.Context, id string) (User, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		// ถ้าแปลงไม่สำเร็จ (เช่น Frontend ส่งคำแปลกๆ ที่ไม่ใช่รูปแบบ UUID มา) ให้ส่ง Error กลับ
		return User{}, fmt.Errorf("invalid uuid format: %w", err)
	}
	return r.q.GetUserByID(ctx, userUUID)
}

func (r *UserRepository) CreateUser(ctx context.Context, id string, email string, fullName *string) (User, error) {
	// 1. แปลง ID ที่เป็น String ให้เป็น UUID
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return User{}, fmt.Errorf("invalid uuid format: %w", err)
	}

	// 2. นำข้อมูลมาจัดเรียงใส่ Struct ที่ sqlc สร้างเตรียมไว้ให้ (CreateUserParams)
	var nullFullName sql.NullString
	if fullName != nil {
		// ถ้า Frontend ส่งชื่อมา (ไม่เป็น null)
		nullFullName.String = *fullName
		nullFullName.Valid = true
	} else {
		// ถ้า Frontend ไม่ได้ส่งชื่อมา (เป็น null)
		nullFullName.String = ""
		nullFullName.Valid = false
	}

	arg := CreateUserParams{
		ID:       userUUID,
		Email:    email,
		FullName: nullFullName,
	}

	// 3. สั่งรันคำสั่ง Insert ผ่าน sqlc
	return r.q.CreateUser(ctx, arg)
}

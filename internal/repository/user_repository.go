package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ไฟล์ที่ sqlc สร้างให้ (models.go, query.sql.go) อยู่ใน package นี้แล้ว
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
		return User{}, fmt.Errorf("invalid uuid format: %w", err)
	}

	// รับค่าจาก sqlc (จะได้เป็น GetUserByIDRow กลับมา)
	row, err := r.q.GetUserByID(ctx, userUUID)
	if err != nil {
		return User{}, err
	}

	return row.User, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, email string, fullName string) (User, error) {
	arg := CreateUserParams{
		Email:    email,
		FullName: fullName,
	}

	// 3. สั่งรันคำสั่ง Insert ผ่าน sqlc
	return r.q.CreateUser(ctx, arg)
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, email string, fullName string) (User, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return User{}, fmt.Errorf("Invalid uuid format: %w", err)
	}

	arg := UpdateUserParams{
		ID:       userUUID,
		Email:    email,
		FullName: fullName,
	}

	return r.q.UpdateUser(ctx, arg)
}

func (r *UserRepository) SoftDeleteUser(ctx context.Context, id string) error {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("Invalid uuid format: %w", err)
	}
	return r.q.SoftDeleteUser(ctx, userUUID)
}

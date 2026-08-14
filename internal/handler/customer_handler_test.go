package handler

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jnarongdech/go-backend-api/internal/dto"
	"github.com/jnarongdech/go-backend-api/internal/repository"
)

// ---------------------------------------------------------
// Create Mock Service
// ---------------------------------------------------------
type MockCustomerService struct {
	MockData repository.Customer
	MockErr  error
}

func (m *MockCustomerService) GetCustomerByID(ctx context.Context, id uuid.UUID) (repository.Customer, error) {
	return m.MockData, m.MockErr
}

// ---------------------------------------------------------
// Unit Tests
// ---------------------------------------------------------

// case 1: Happy Path (success)
func TestGetCustomerByID_Success(t *testing.T) {
	// จัดฉาก: ให้ Mock คืนค่าสำเร็จ
	mockSvc := &MockCustomerService{
		MockData: repository.Customer{
			ID:    uuid.New(),
			Name:  "สมชาย ใจดี",
			Email: "somchai@example.com",
			CompanyName: sql.NullString{
				String: "STEEL-FACTORY",
				Valid:  true,
			},
		},
		MockErr: nil,
	}

	handler := NewCustomerHandler(mockSvc)
	app := fiber.New()
	app.Get("/api/v1/customers/:id", handler.GetCustomerByID)

	// ยิง Request
	req := httptest.NewRequest("GET", "/api/v1/customers/1", nil)
	resp, _ := app.Test(req)

	// 1. เช็ค Status Code ว่าต้องเป็น 200
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("คาดหวัง Status 200 แต่ได้ %d", resp.StatusCode)
	}

	// 2. (โบนัส) เช็ค Body ว่ามีคำว่า "Success" ตามที่เราเขียนไว้ไหม
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Success!") {
		t.Errorf("คาดหวังให้มีคำว่า Success! แต่ได้ body: %s", string(body))
	}
}

// เคสที่ 2: Unhappy Path (not found customer หรือ DB Broken)
func TestGetCustomerByID_NotFound(t *testing.T) {
	// จัดฉาก: ให้ Mock โยน Error กลับมา
	mockSvc := &MockCustomerService{
		MockData: repository.Customer{},
		MockErr:  errors.New("database connection failed"),
	}

	handler := NewCustomerHandler(mockSvc)
	app := fiber.New()
	app.Get("/api/v1/customers/:id", handler.GetCustomerByID)

	req := httptest.NewRequest("GET", "/api/v1/customers/99", nil)
	resp, _ := app.Test(req)

	// เช็ค Status Code ว่าต้องเป็น 400 (ตามตรรกะใน Handler ของคุณ)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("คาดหวัง Status 400 แต่ได้ %d", resp.StatusCode)
	}

	// เช็ค Body ว่าแจ้งเตือนถูกต้องไหม
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Customer not found") {
		t.Errorf("คาดหวังให้มีคำว่า Customer not found แต่ได้ body: %s", string(body))
	}
}

// เคสที่ 3: Edge Case (ไม่ส่ง ID มา)
func TestGetCustomerByID_MissingID(t *testing.T) {
	mockSvc := &MockCustomerService{
		MockData: repository.Customer{},
		MockErr:  nil,
	}

	handler := NewCustomerHandler(mockSvc)
	app := fiber.New()

	// ตั้งใจสร้าง Route ที่รับ ID หรือไม่รับก็ได้ (ใส่ ? ต่อท้าย)
	// เพื่อหลอกให้ Fiber ปล่อยผ่านเข้าไปถึงบรรทัด if customerID == "" ใน Handler
	app.Get("/api/v1/customers/:id?", handler.GetCustomerByID)

	// ยิง Request แบบไม่ส่ง ID
	req := httptest.NewRequest("GET", "/api/v1/customers/", nil)
	resp, _ := app.Test(req)

	// เช็ค Status Code ว่าต้องเป็น 400
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("คาดหวัง Status 400 แต่ได้ %d", resp.StatusCode)
	}

	// เช็ค Body
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Customer ID is required") {
		t.Errorf("คาดหวังให้มีคำว่า Customer ID is required แต่ได้ body: %s", string(body))
	}
}

func (m *MockCustomerService) CreateCustomer(ctx context.Context, req dto.CreateCustomerRequest) (repository.Customer, error) {
	return m.MockData, m.MockErr
}

func (m *MockCustomerService) UpdateCustomer(ctx context.Context, req dto.UpdateCustomerRequest) error {
	return nil
}

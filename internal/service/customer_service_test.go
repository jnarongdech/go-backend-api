package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jnarongdech/go-backend-api/internal/repository"
)

// Mock Repo
type MockCustomerRepository struct {
	MockData repository.Customer
	MockErr  error
}

func (m *MockCustomerRepository) GetCustomerByID(ctx context.Context, id uuid.UUID) (repository.Customer, error) {
	return m.MockData, m.MockErr
}

// Start Unit Test : GetCustomerByID
// case 1 : Happy Path
func TestGetCustomerByID_Success(t *testing.T) {
	mockID := uuid.New()
	mockRepo := &MockCustomerRepository{
		MockData: repository.Customer{
			ID:   mockID,
			Name: "สมชาย ใจดี",
		},
		MockErr: nil,
	}

	service := NewCustomerService(mockRepo)

	// try to call service
	customer, err := service.GetCustomerByID(context.Background(), mockID)
	if err != nil {
		t.Fatalf("No error should occur: %v", err)
	}
	if customer.Name != "สมชาย ใจดี" {
		t.Errorf("Expected name: สมชาย ใจดี but got %s", customer.Name)
	}
}

// case 2: id cannot be empty
func TestGetCustomerByID_Empty(t *testing.T) {
	mockRepo := &MockCustomerRepository{}
	service := NewCustomerService(mockRepo)

	_, err := service.GetCustomerByID(context.Background(), uuid.Nil)

	if err == nil {
		t.Fatalf("Expected an error but none occurred.")
	}
	if err.Error() != "Customer ID cannot be empty." {
		t.Errorf("The error does not match.")
	}
}

// case 3: Not found customer (condition: sql.ErrNoRows)
func TestGetCustomerByID_NotFound(t *testing.T) {
	mockRepo := &MockCustomerRepository{
		MockData: repository.Customer{},
		MockErr:  sql.ErrNoRows,
	}
	service := NewCustomerService(mockRepo)
	_, err := service.GetCustomerByID(context.Background(), uuid.New())

	if err == nil {
		t.Fatalf("Expected an error but none occurred.")
	}

	if err.Error() != "Not found customer" {
		t.Errorf("The message error do not match, but got: %s", err.Error())
	}
}

// case 4: DB Error (General Error)
func TestGetCustomerByID_DatabaseError(t *testing.T) {
	mockRepo := &MockCustomerRepository{
		MockData: repository.Customer{},
		MockErr:  errors.New("connection timeout"),
	}

	service := NewCustomerService(mockRepo)
	_, err := service.GetCustomerByID(context.Background(), uuid.New())

	if err == nil {
		t.Fatalf("Expected an error but none occurred.")
	}

	if err.Error() != "An error has occurred, please try again." {
		t.Errorf("The message error do not match, but got: %s", err.Error())
	}
}

// end of unit test: GetCustomerByID

func (m *MockCustomerRepository) CreateCustomer(ctx context.Context, arg repository.CreateCustomerParams) (repository.Customer, error) {
	return m.MockData, m.MockErr
}

func (m *MockCustomerRepository) UpdateCustomer(context.Context, repository.UpdateCustomerParams) error {
	return m.MockErr
}

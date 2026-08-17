package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jnarongdech/go-backend-api/internal/dto"
	"github.com/jnarongdech/go-backend-api/internal/repository"
	constants "github.com/jnarongdech/go-backend-api/pkg/consts"
	"github.com/jnarongdech/go-backend-api/pkg/errs"
	"github.com/jnarongdech/go-backend-api/pkg/utils"
)

// 1. Repository Interface
type CustomerRepository interface {
	GetCustomerByID(ctx context.Context, id uuid.UUID) (repository.Customer, error)
	CreateCustomer(ctx context.Context, arg repository.CreateCustomerParams) (repository.Customer, error)
	UpdateCustomer(ctx context.Context, arg repository.UpdateCustomerParams) error
	ExecTx(ctx context.Context, fn func(*repository.Queries) error) error
}

// 2. Service Interface (เพื่อให้ Handler เรียกใช้)
type CustomerService interface {
	GetCustomerByID(ctx context.Context, id uuid.UUID) (repository.Customer, error)
	CreateCustomer(ctx context.Context, req dto.CreateCustomerRequest) (repository.Customer, error)
	UpdateCustomer(ctx context.Context, req dto.UpdateCustomerRequest) error
}

// 3. Struct ตัวพิมพ์เล็ก (Private Struct)
type customerService struct {
	store CustomerRepository
}

// 4. Constructor คืนค่าเป็น Interface
func NewCustomerService(store CustomerRepository) CustomerService {
	return &customerService{store: store}
}

func (s *customerService) GetCustomerByID(ctx context.Context, id uuid.UUID) (repository.Customer, error) {

	customer, err := s.store.GetCustomerByID(ctx, id)
	if err != nil {
		// log for dev check
		log.Printf("[ERROR] GetCustomerByID failed for id %s: %v", id, err)

		if errors.Is(err, sql.ErrNoRows) {
			return repository.Customer{}, errs.NewNotFound(constants.ErrResourceNotFound, err)
		}

		return repository.Customer{}, errs.NewInternal(constants.ErrInternalServer, err)
	}

	return customer, nil
}

func (s *customerService) CreateCustomer(ctx context.Context, req dto.CreateCustomerRequest) (repository.Customer, error) {
	arg := repository.CreateCustomerParams{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       utils.ValueToNullString(req.Phone),
		CompanyName: utils.ValueToNullString(req.CompanyName),
		Address:     utils.ValueToNullString(req.Address),
		City:        utils.ValueToNullString(req.City),
		PostalCode:  utils.ValueToNullString(req.PostalCode),
		Country:     utils.ValueToNullString(req.Country),
	}

	// Pass the (converted) arg to the Repository for handling.
	customer, err := s.store.CreateCustomer(ctx, arg)
	if err != nil {
		log.Printf("[ERROR] CreateCustomer failed: %v", err)

		// You could add a check for duplicate emails here.
		return repository.Customer{}, errs.NewInternal(constants.ErrInternalServer, err)
	}

	return customer, nil
}

func (s *customerService) UpdateCustomer(ctx context.Context, req dto.UpdateCustomerRequest) error {
	arg := repository.UpdateCustomerParams{
		ID:          req.ID,
		Name:        req.Name,
		Email:       req.Email,
		Phone:       utils.ValueToNullString(req.Phone),
		CompanyName: utils.ValueToNullString(req.CompanyName),
		Address:     utils.ValueToNullString(req.Address),
		City:        utils.ValueToNullString(req.City),
		PostalCode:  utils.ValueToNullString(req.PostalCode),
		Country:     utils.ValueToNullString(req.Country),
	}

	err := s.store.UpdateCustomer(ctx, arg)
	if err != nil {
		log.Printf("[ERROR] UpdateCustomer failed for ID %s: %v", req.ID, err)

		// customer not found
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewNotFound(constants.ErrCustomerNotFound, err)
		}

		// duplicate email
		// if strings.Contains(err.Error(), "duplicate key value") {
		// 	return errors.New("This email is already taken.")
		// }

		return errs.NewInternal(constants.ErrInternalServer, err)
	}

	return nil
}

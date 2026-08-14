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
)

type ICustomerRepository interface {
	GetCustomerByID(ctx context.Context, id uuid.UUID) (repository.Customer, error)
	CreateCustomer(ctx context.Context, arg repository.CreateCustomerParams) (repository.Customer, error)
	UpdateCustomer(ctx context.Context, arg repository.UpdateCustomerParams) error
}

type CustomerService struct {
	store ICustomerRepository
}

// Constructor
func NewCustomerService(store ICustomerRepository) *CustomerService {
	return &CustomerService{store: store}
}

func (s *CustomerService) GetCustomerByID(ctx context.Context, id uuid.UUID) (repository.Customer, error) {
	if id == uuid.Nil {
		return repository.Customer{}, errors.New(constants.ErrCustomerIDEmpty)
	}

	customer, err := s.store.GetCustomerByID(ctx, id)
	if err != nil {

		// log for dev check
		log.Printf("[ERROR] GetCustomerByID failed for id %s: %v", id, err)

		if errors.Is(err, sql.ErrNoRows) {
			// ถ้าใช่ ให้ส่งข้อความที่เข้าใจง่ายกลับไป
			return repository.Customer{}, errors.New("Not found customer")
		}

		return repository.Customer{}, errors.New("An error has occurred, please try again.")
	}

	return customer, nil
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req dto.CreateCustomerRequest) (repository.Customer, error) {
	if req.Name == "" {
		return repository.Customer{}, errors.New("Customer name cannot be empty.")
	}
	if req.Email == "" {
		return repository.Customer{}, errors.New("Customer email cannot be empty.")
	}

	// Convert the DTO into Repository parameters.
	arg := repository.CreateCustomerParams{
		Name:  req.Name,
		Email: req.Email,
		Phone: sql.NullString{
			String: req.Phone,
			Valid:  req.Phone != "",
		},
		CompanyName: sql.NullString{
			String: req.CompanyName,
			Valid:  req.CompanyName != "",
		},
		Address: sql.NullString{
			String: req.Address,
			Valid:  req.Address != "",
		},
		City: sql.NullString{
			String: req.City,
			Valid:  req.City != "",
		},
		PostalCode: sql.NullString{
			String: req.PostalCode,
			Valid:  req.PostalCode != "",
		},
		Country: sql.NullString{
			String: req.Country,
			Valid:  req.Country != "",
		},
	}

	// Pass the (converted) arg to the Repository for handling.
	customer, err := s.store.CreateCustomer(ctx, arg)
	if err != nil {
		log.Printf("[ERROR] CreateCustomer failed: %v", err)

		// You could add a check for duplicate emails here.
		return repository.Customer{}, errors.New("Unable to create customer data.")
	}

	return customer, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, req dto.UpdateCustomerRequest) error {
	if req.ID == uuid.Nil {
		return errors.New(constants.ErrCustomerIDEmpty)
	}
	if req.Name == "" {
		return errors.New(constants.ErrCustomerNameEmpty)
	}
	if req.Email == "" {
		return errors.New(constants.ErrCustomerEmailEmpty)
	}

	arg := repository.UpdateCustomerParams{
		ID:    req.ID,
		Name:  req.Name,
		Email: req.Email,
		Phone: sql.NullString{
			String: req.Phone,
			Valid:  req.Phone != "",
		},
		CompanyName: sql.NullString{
			String: req.CompanyName,
			Valid:  req.CompanyName != "",
		},
		Address: sql.NullString{
			String: req.Address,
			Valid:  req.Address != "",
		},
		City: sql.NullString{
			String: req.City,
			Valid:  req.City != "",
		},
		PostalCode: sql.NullString{
			String: req.PostalCode,
			Valid:  req.PostalCode != "",
		},
		Country: sql.NullString{
			String: req.Country,
			Valid:  req.Country != "",
		},
	}

	err := s.store.UpdateCustomer(ctx, arg)
	if err != nil {
		log.Printf("[ERROR] UpdateCustomer failed for ID %s: %v", req.ID, err)

		// customer not found
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(constants.ErrCustomerNotFound)
		}

		// duplicate email
		// if strings.Contains(err.Error(), "duplicate key value") {
		// 	return errors.New("This email is already taken.")
		// }

		return errors.New(constants.ErrUnableUpdateData)
	}

	return nil
}

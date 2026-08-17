package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jnarongdech/go-backend-api/internal/dto"
	"github.com/jnarongdech/go-backend-api/internal/repository"
	constants "github.com/jnarongdech/go-backend-api/pkg/consts"
	"github.com/jnarongdech/go-backend-api/pkg/errs"
	"github.com/jnarongdech/go-backend-api/pkg/utils"
)

type ProductService struct {
	store *repository.Store
}

func NewProductService(store *repository.Store) *ProductService {
	return &ProductService{
		store: store, // จับยัดใส่ struct
	}
}

func (s *ProductService) GetProducts(ctx context.Context) ([]repository.Product, error) {
	products, err := s.store.GetProducts(ctx)

	if err != nil {
		log.Printf("[ERROR] GetProducts failed: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return []repository.Product{}, errs.NewNotFound(constants.ErrResourceNotFound, err)
		}

		return []repository.Product{}, errs.NewInternal(constants.ErrInternalServer, err)
	}

	if products == nil {
		return []repository.Product{}, nil
	}

	return products, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error) {
	if s.store == nil {
		fmt.Println("s.store เป็น nil (ลืม Inject Database แน่ๆ)")
	}

	arg := repository.CreateProductParams{
		Name:                req.Name,
		Description:         utils.PointerToNullString(req.Description),
		Category:            utils.PointerToNullString(req.Category),
		BasePrice:           req.BasePrice,
		IsCustomizable:      utils.PointerToNullBool(req.IsCustomizable),
		CustomizationFields: utils.JsonToNullRawMessage(req.CustomizationFields),
	}

	row, err := s.store.CreateProduct(ctx, arg)
	if err != nil {
		log.Printf("\n[ERROR] Create product failed: %v", err)
		return dto.ProductResponse{}, errs.NewInternal(constants.ErrInternalServer, err)
	}

	var customFields json.RawMessage
	if row.CustomizationFields.Valid {
		customFields = row.CustomizationFields.RawMessage
	} else {
		customFields = nil // ถ้าไม่มีข้อมูล ให้พ่น null
	}

	result := dto.ProductResponse{
		ID:                  row.ID.String(),
		Name:                row.Name,
		Description:         utils.NullStringToPointer(row.Description),
		Category:            utils.NullStringToPointer(row.Category),
		BasePrice:           row.BasePrice,
		IsCustomizable:      utils.NullBoolToPointer(row.IsCustomizable),
		CustomizationFields: customFields,
	}

	return result, nil
}

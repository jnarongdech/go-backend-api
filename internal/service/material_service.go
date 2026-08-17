package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jnarongdech/go-backend-api/internal/dto"
	"github.com/jnarongdech/go-backend-api/internal/repository"
	constants "github.com/jnarongdech/go-backend-api/pkg/consts"
	"github.com/jnarongdech/go-backend-api/pkg/errs"
	"github.com/jnarongdech/go-backend-api/pkg/utils"
)

// define repo interface
type MaterialRepository interface {
	CreateMaterial(ctx context.Context, arg repository.CreateMaterialParams) (repository.Material, error)
	GetMaterials(ctx context.Context) ([]repository.Material, error)
	GetMaterialByID(ctx context.Context, id uuid.UUID) (repository.Material, error)
	ExecTx(ctx context.Context, fn func(*repository.Queries) error) error
}

// define service interface
type MaterialService interface {
	CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest) (*dto.MaterialResponse, error)
	GetMaterails(ctx context.Context) ([]dto.MaterialResponse, error)
	GetMaterialByID(ctx context.Context, id uuid.UUID) (dto.MaterialResponse, error)
}

// define private struct
type materialService struct {
	store MaterialRepository
}

func NewMaterialService(store MaterialRepository) MaterialService {
	return &materialService{
		store: store,
	}
}

func (s *materialService) CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest) (*dto.MaterialResponse, error) {
	arg := repository.CreateMaterialParams{
		Name:           req.Name,
		ThicknessMm:    utils.FloatPointerToNullString(&req.ThicknessMM),
		Grade:          utils.ValueToNullString(req.Grade),
		Description:    utils.PointerToNullString(&req.Description),
		CostPerKg:      utils.FloatValueToNullString(req.CostPerKg),
		StockQtyKg:     utils.FloatPointerToNullString(&req.StockQtyKg),
		ReorderLevelKg: utils.FloatPointerToNullString(&req.ReorderLevelKg),
	}

	row, err := s.store.CreateMaterial(ctx, arg)
	if err != nil {
		return &dto.MaterialResponse{}, errs.NewInternal(constants.ErrInternalServer, err)
	}

	result := dto.MaterialResponse{
		ID:             row.ID,
		Name:           row.Name,
		ThicknessMM:    req.ThicknessMM, // ปลอดภัยกว่า
		Grade:          req.Grade,       // ปลอดภัยกว่า
		Description:    req.Description,
		CostPerKg:      req.CostPerKg,
		StockQtyKg:     req.StockQtyKg,
		ReorderLevelKg: req.ReorderLevelKg,
	}

	return &result, nil
}

func (s *materialService) GetMaterails(ctx context.Context) ([]dto.MaterialResponse, error) {
	data, err := s.store.GetMaterials(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []dto.MaterialResponse{}, errs.NewNotFound(constants.ErrResourceNotFound, err)
		}
	}

	materialList := make([]dto.MaterialResponse, 0, len(data))
	for _, item := range data {
		mappedItem := dto.MaterialResponse{
			ID:             item.ID,
			Name:           item.Name,
			ThicknessMM:    *utils.NullStringToFloatPointer(item.ThicknessMm),
			Grade:          *utils.NullStringToPointer(item.Grade),
			Description:    *utils.NullStringToPointer(item.Description),
			CostPerKg:      *utils.NullStringToFloatPointer(item.CostPerKg),
			StockQtyKg:     *utils.NullStringToFloatPointer(item.StockQtyKg),
			ReorderLevelKg: *utils.NullStringToFloatPointer(item.ReorderLevelKg),
			CreatedAt:      *utils.NullTimeToPointer(item.CreatedAt),
			UpdatedAt:      *utils.NullTimeToPointer(item.UpdatedAt),
		}
		materialList = append(materialList, mappedItem)
	}

	return materialList, nil
}

func (s *materialService) GetMaterialByID(ctx context.Context, id uuid.UUID) (dto.MaterialResponse, error) {
	item, err := s.store.GetMaterialByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.MaterialResponse{}, errs.NewNotFound(constants.ErrResourceNotFound, err)
		}
		return dto.MaterialResponse{}, errs.NewInternal(constants.ErrInternalServer, err)
	}

	result := dto.MaterialResponse{
		ID:             item.ID,
		Name:           item.Name,
		ThicknessMM:    *utils.NullStringToFloatPointer(item.ThicknessMm),
		Grade:          *utils.NullStringToPointer(item.Grade),
		Description:    *utils.NullStringToPointer(item.Description),
		CostPerKg:      *utils.NullStringToFloatPointer(item.CostPerKg),
		StockQtyKg:     *utils.NullStringToFloatPointer(item.StockQtyKg),
		ReorderLevelKg: *utils.NullStringToFloatPointer(item.ReorderLevelKg),
		CreatedAt:      *utils.NullTimeToPointer(item.CreatedAt),
		UpdatedAt:      *utils.NullTimeToPointer(item.UpdatedAt),
	}

	return result, nil
}

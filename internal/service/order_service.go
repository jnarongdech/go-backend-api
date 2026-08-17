package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jnarongdech/go-backend-api/internal/dto"
	"github.com/jnarongdech/go-backend-api/internal/repository"
	constants "github.com/jnarongdech/go-backend-api/pkg/consts"
	"github.com/jnarongdech/go-backend-api/pkg/errs"
	"github.com/jnarongdech/go-backend-api/pkg/utils"
)

type OrderService struct {
	store *repository.Store
}

func NewOrderService(store *repository.Store) *OrderService {
	return &OrderService{store: store}
}

func (s *OrderService) GetOrders(ctx context.Context) ([]dto.OrderResponse, error) {
	data, err := s.store.GetOrders(ctx)
	if err != nil {
		log.Printf("[ERROR] GetOrders failed: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return []dto.OrderResponse{}, errs.NewNotFound(constants.ErrResourceNotFound, err)
		}
		return []dto.OrderResponse{}, errs.NewInternal(constants.ErrInternalServer, err)
	}

	result := make([]dto.OrderResponse, 0, len(data))
	for _, dataItem := range data {
		// เตรียมข้อมูล (ป้องกัน nil pointer ให้ครบทุกตัว)
		var resTotalPrice *float64
		if dataItem.TotalPrice.Valid {
			val, err := strconv.ParseFloat(dataItem.TotalPrice.String, 64)
			if err == nil {
				resTotalPrice = &val
			}
		}

		var resNotes *string
		if dataItem.Notes.Valid {
			notesVal := dataItem.Notes.String
			resNotes = &notesVal
		}

		var resOrderDate *time.Time
		if dataItem.OrderDate.Valid {
			orderDateVal := dataItem.OrderDate.Time
			resOrderDate = &orderDateVal
		}

		var resExpectedDate *time.Time
		if dataItem.ExpectedCompletionDate.Valid {
			exDateVal := dataItem.ExpectedCompletionDate.Time
			resExpectedDate = &exDateVal
		}

		var resActualDate *time.Time
		if dataItem.ActualCompletionDate.Valid {
			acDateVal := dataItem.ActualCompletionDate.Time
			resActualDate = &acDateVal
		}

		mappedItem := dto.OrderResponse{
			ID:                     dataItem.ID,
			OrderNumber:            dataItem.OrderNumber,
			CustomerID:             dataItem.CustomerID,
			OrderType:              dataItem.OrderType,
			Status:                 dataItem.Status,
			TotalPrice:             resTotalPrice,
			Notes:                  resNotes,
			OrderDate:              resOrderDate,
			ExpectedCompletionDate: resExpectedDate,
			ActualCompletionDate:   resActualDate,
		}

		result = append(result, mappedItem)
	}

	return result, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id uuid.UUID) (dto.OrderResponse, error) {
	row, err := s.store.GetOrderWithItems(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("[ERROR] Get order by id failed: %v", err)
			return dto.OrderResponse{}, errs.NewNotFound(constants.ErrResourceNotFound, err) // 404
		}
	}

	var orderItems []dto.OrderItemsResponse
	if len(row.Items) > 0 {
		err = json.Unmarshal(row.Items, &orderItems)
		if err != nil {
			return dto.OrderResponse{}, errs.NewInternal("Unable to read item list data", err)
		}
	}

	result := dto.OrderResponse{
		ID:                     row.ID,
		OrderNumber:            row.OrderNumber,
		CustomerID:             row.CustomerID,
		OrderType:              row.OrderType,
		Status:                 row.Status,
		TotalPrice:             utils.NullStringToFloatPointer(row.TotalPrice),
		Notes:                  utils.NullStringToPointer(row.Notes),
		OrderDate:              utils.NullTimeToPointer(row.OrderDate),
		ExpectedCompletionDate: utils.NullTimeToPointer(row.ExpectedCompletionDate),
		ActualCompletionDate:   utils.NullTimeToPointer(row.ActualCompletionDate),
		Items:                  orderItems,
	}

	return result, nil
}

func (s *OrderService) CreateOrderWithItems(ctx context.Context, req dto.CreateOrderRequest) error {
	err := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		// เตรียมข้อมูล (ป้องกัน nil pointer ให้ครบทุกตัว)
		var dbTotalPrice sql.NullString
		if req.TotalPrice != nil {
			dbTotalPrice = sql.NullString{
				String: fmt.Sprintf("%.2f", *req.TotalPrice),
				Valid:  true,
			}
		}

		var dbNotes sql.NullString
		if req.Notes != nil {
			dbNotes = sql.NullString{String: *req.Notes, Valid: true}
		}

		var dbOrderDate sql.NullTime
		if req.OrderDate != nil {
			dbOrderDate = sql.NullTime{Time: *req.OrderDate, Valid: true}
		}

		var dbExpectedDate sql.NullTime
		if req.ExpectedCompletionDate != nil {
			dbExpectedDate = sql.NullTime{Time: *req.ExpectedCompletionDate, Valid: true}
		}

		var dbActualDate sql.NullTime
		if req.ActualCompletionDate != nil {
			dbActualDate = sql.NullTime{Time: *req.ActualCompletionDate, Valid: true}
		}

		// time.Now().Format("0601") จะได้ปี 2 หลัก + เดือน 2 หลัก (เช่น 2608)
		currentMonth := time.Now().Format("0601")
		randomCode := generateRandomCode(6)

		// จะได้ผลลัพธ์ประมาณ: "ORD-2608-X7K9"
		orderNumber := fmt.Sprintf("ORD-%s-%s", currentMonth, randomCode)

		orderArg := repository.CreateOrderParams{
			CustomerID:             req.CustomerID,
			OrderNumber:            orderNumber,
			OrderType:              req.OrderType,
			Status:                 req.Status,
			TotalPrice:             dbTotalPrice,
			Notes:                  dbNotes,
			OrderDate:              dbOrderDate,
			ExpectedCompletionDate: dbExpectedDate,
			ActualCompletionDate:   dbActualDate,
		}

		createOrder, errCreateOrder := qtx.CreateOrder(ctx, orderArg)
		if errCreateOrder != nil {
			log.Printf("[ERROR] CreateOrder failed: %v", errCreateOrder)
			return errs.NewInternal(constants.ErrInternalServer, errCreateOrder)
		}

		for _, itemReq := range req.Items {
			itemArg := repository.CreateOrderItemsParams{
				OrderID:      createOrder.ID,
				ProductID:    itemReq.ProductID,
				Quantity:     itemReq.Quantity,
				PricePerUnit: fmt.Sprintf("%.2f", itemReq.PricePerUnit),
			}

			_, errItem := qtx.CreateOrderItems(ctx, itemArg)
			if errItem != nil {
				return errs.NewInternal(constants.ErrInternalServer, errItem)
			}
		}

		return nil
	})

	// เช็ค Error จาก Transaction
	if err != nil {
		return errs.NewInternal(constants.ErrInternalServer, err)
	}

	return nil
}

func (s *OrderService) UpdateOrderWithItems(ctx context.Context, orderID uuid.UUID, req dto.UpdateOrderRequest) error {
	err := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		// เตรียมข้อมูล (ป้องกัน nil pointer ให้ครบทุกตัว)
		var dbTotalPrice sql.NullString
		if req.TotalPrice != nil {
			dbTotalPrice = sql.NullString{
				String: fmt.Sprintf("%.2f", *req.TotalPrice),
				Valid:  true,
			}
		}

		var dbNotes sql.NullString
		if req.Notes != nil {
			dbNotes = sql.NullString{String: *req.Notes, Valid: true}
		}

		var dbOrderDate sql.NullTime
		if req.OrderDate != nil {
			dbOrderDate = sql.NullTime{Time: *req.OrderDate, Valid: true}
		}

		var dbExpectedDate sql.NullTime
		if req.ExpectedCompletionDate != nil {
			dbExpectedDate = sql.NullTime{Time: *req.ExpectedCompletionDate, Valid: true}
		}

		var dbActualDate sql.NullTime
		if req.ActualCompletionDate != nil {
			dbActualDate = sql.NullTime{Time: *req.ActualCompletionDate, Valid: true}
		}

		orderArg := repository.UpdateOrderParams{
			ID:                     orderID,
			CustomerID:             req.CustomerID,
			OrderType:              req.OrderType,
			Status:                 req.Status,
			TotalPrice:             dbTotalPrice,
			Notes:                  dbNotes,
			OrderDate:              dbOrderDate,
			ExpectedCompletionDate: dbExpectedDate,
			ActualCompletionDate:   dbActualDate,
		}

		err := qtx.UpdateOrder(ctx, orderArg)
		if err != nil {
			log.Printf("[ERROR] Update order failed for ID %s: %v", req.ID, err)
			return errs.NewInternal(constants.ErrInternalServer, err)
		}

		// delete items and add items
		errDeleteOrderItems := qtx.DeleteOrderItemsByOrderID(ctx, orderID)
		if errDeleteOrderItems != nil {
			return errors.New("Unable to delete order items.")
		}

		for _, itemReq := range req.Items {
			itemArg := repository.CreateOrderItemsParams{
				OrderID:      orderID,
				ProductID:    itemReq.ProductID,
				Quantity:     itemReq.Quantity,
				PricePerUnit: fmt.Sprintf("%.2f", itemReq.PricePerUnit),
			}

			_, err := qtx.CreateOrderItems(ctx, itemArg)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func generateRandomCode(length int) string {
	const charset = "ABCDEFGHJKMNPQRSTUVWXYZ123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

package repository

import (
	"context"
	"fmt"

	"github.com/IainMosima/gomart/domains/order/entity"
	domainRepo "github.com/IainMosima/gomart/domains/order/repository"
	db "github.com/IainMosima/gomart/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrderRepositoryImpl struct {
	store db.Store
}

func NewOrderRepository(store db.Store) domainRepo.OrderRepository {
	return &OrderRepositoryImpl{
		store: store,
	}
}

func (r *OrderRepositoryImpl) CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	var status pgtype.Text
	status = pgtype.Text{
		String: string(order.Status),
		Valid:  true,
	}

	totalAmount := pgtype.Numeric{}
	totalAmountStr := fmt.Sprintf("%.2f", order.TotalAmount)
	if err := totalAmount.Scan(totalAmountStr); err != nil {
		return nil, err
	}

	params := db.CreateOrderParams{
		CustomerID:  order.CustomerID,
		OrderNumber: order.OrderNumber,
		Status:      status,
		TotalAmount: totalAmount,
	}

	result, err := r.store.CreateOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.convertOrderToEntity(result), nil
}

func (r *OrderRepositoryImpl) CreateOrderItem(ctx context.Context, orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	unitPrice := pgtype.Numeric{}
	unitPriceStr := fmt.Sprintf("%.2f", orderItem.UnitPrice)
	if err := unitPrice.Scan(unitPriceStr); err != nil {
		return nil, err
	}

	totalPrice := pgtype.Numeric{}
	totalPriceStr := fmt.Sprintf("%.2f", orderItem.TotalPrice)
	if err := totalPrice.Scan(totalPriceStr); err != nil {
		return nil, err
	}

	params := db.CreateOrderItemParams{
		OrderID:    orderItem.OrderID,
		ProductID:  orderItem.ProductID,
		Quantity:   orderItem.Quantity,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
	}

	result, err := r.store.CreateOrderItem(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.convertOrderItemToEntity(result), nil
}

func (r *OrderRepositoryImpl) GetOrder(ctx context.Context, orderID uuid.UUID) (*entity.Order, error) {
	result, err := r.store.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return r.convertOrderToEntity(result), nil
}

func (r *OrderRepositoryImpl) convertOrderToEntity(dbOrder db.Order) *entity.Order {
	order := &entity.Order{
		OrderID:     dbOrder.OrderID,
		CustomerID:  dbOrder.CustomerID,
		OrderNumber: dbOrder.OrderNumber,
		CreatedAt:   dbOrder.CreatedAt.Time,
		IsDeleted:   dbOrder.IsDeleted.Bool,
	}

	if dbOrder.Status.Valid {
		order.Status = entity.OrderStatus(dbOrder.Status.String)
	}

	if dbOrder.TotalAmount.Valid {
		totalAmount, err := dbOrder.TotalAmount.Float64Value()
		if err == nil {
			order.TotalAmount = totalAmount.Float64
		}
	}

	if dbOrder.UpdatedAt.Valid {
		updatedAt := dbOrder.UpdatedAt.Time
		order.UpdatedAt = &updatedAt
	}

	return order
}

func (r *OrderRepositoryImpl) convertOrderItemToEntity(dbOrderItem db.OrderItem) *entity.OrderItem {
	orderItem := &entity.OrderItem{
		OrderItemID: dbOrderItem.OrderItemID,
		OrderID:     dbOrderItem.OrderID,
		ProductID:   dbOrderItem.ProductID,
		Quantity:    dbOrderItem.Quantity,
		CreatedAt:   dbOrderItem.CreatedAt.Time,
		IsDeleted:   dbOrderItem.IsDeleted.Bool,
	}

	if dbOrderItem.UnitPrice.Valid {
		unitPrice, err := dbOrderItem.UnitPrice.Float64Value()
		if err == nil {
			orderItem.UnitPrice = unitPrice.Float64
		}
	}

	if dbOrderItem.TotalPrice.Valid {
		totalPrice, err := dbOrderItem.TotalPrice.Float64Value()
		if err == nil {
			orderItem.TotalPrice = totalPrice.Float64
		}
	}

	return orderItem
}

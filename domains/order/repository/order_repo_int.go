package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/order/entity"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	CreateOrderItem(ctx context.Context, orderItem *entity.OrderItem) (*entity.OrderItem, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*entity.Order, error)
}

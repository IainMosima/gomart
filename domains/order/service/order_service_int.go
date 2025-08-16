package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/order/schema"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *schema.CreateOrderRequest) (*schema.OrderResponse, error)
	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (*schema.OrderStatusResponse, error)
}

package order

import (
	"context"

	"github.com/IainMosima/gomart/domains/order/repository"
	"github.com/IainMosima/gomart/domains/order/schema"
	"github.com/google/uuid"
)

type OrderServiceImpl struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepository repository.OrderRepository) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo: orderRepository,
	}
}

func (o OrderServiceImpl) CreateOrder(ctx context.Context, req *schema.CreateOrderRequest) (*schema.OrderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderServiceImpl) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (*schema.OrderStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

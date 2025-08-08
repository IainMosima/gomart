package repository

import (
	"context"
	"time"

	"github.com/IainMosima/gomart/domains/order/entity"
	"github.com/google/uuid"
)

type OrderRepository interface {
	// Order CRUD operations
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*entity.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*entity.Order, error)
	UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status entity.OrderStatus) (*entity.Order, error)
	UpdateOrderTotal(ctx context.Context, orderID uuid.UUID, totalAmount float64) (*entity.Order, error)
	SoftDeleteOrder(ctx context.Context, orderID uuid.UUID) error

	// Order listing and filtering
	ListOrders(ctx context.Context) ([]*entity.Order, error)
	ListOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]*entity.Order, error)
	ListOrdersByStatus(ctx context.Context, status entity.OrderStatus) ([]*entity.Order, error)
	ListOrdersByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.Order, error)

	// Order details with related data
	GetOrderWithDetails(ctx context.Context, orderID uuid.UUID) (*entity.OrderWithDetails, error)
	GetOrderSummary(ctx context.Context, orderID uuid.UUID) (*entity.OrderSummary, error)

	// Order item operations
	CreateOrderItem(ctx context.Context, orderItem *entity.OrderItem) (*entity.OrderItem, error)
	GetOrderItem(ctx context.Context, orderItemID uuid.UUID) (*entity.OrderItem, error)
	ListOrderItemsByOrder(ctx context.Context, orderID uuid.UUID) ([]*entity.OrderItem, error)
	ListOrderItemsByProduct(ctx context.Context, productID uuid.UUID) ([]*entity.OrderItem, error)
	GetOrderItemsWithDetails(ctx context.Context, orderID uuid.UUID) ([]*entity.OrderItemWithDetails, error)
	UpdateOrderItem(ctx context.Context, orderItem *entity.OrderItem) (*entity.OrderItem, error)
	UpdateOrderItemQuantity(ctx context.Context, orderItemID uuid.UUID, quantity int32) (*entity.OrderItem, error)
	SoftDeleteOrderItem(ctx context.Context, orderItemID uuid.UUID) error
	SoftDeleteOrderItemsByOrder(ctx context.Context, orderID uuid.UUID) error

	// Order calculations
	GetOrderTotalFromItems(ctx context.Context, orderID uuid.UUID) (float64, error)
	RecalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (*entity.Order, error)

	// Statistics and counting
	CountOrders(ctx context.Context) (int64, error)
	CountOrdersByStatus(ctx context.Context, status entity.OrderStatus) (int64, error)
	CountOrdersByCustomer(ctx context.Context, customerID uuid.UUID) (int64, error)
	CountOrderItems(ctx context.Context) (int64, error)
	CountOrderItemsByOrder(ctx context.Context, orderID uuid.UUID) (int64, error)
	GetOrderStats(ctx context.Context) (*entity.OrderStats, error)

	// Order number generation
	GenerateOrderNumber(ctx context.Context) (string, error)

	// Validation helpers
	OrderExists(ctx context.Context, orderID uuid.UUID) (bool, error)
	OrderBelongsToCustomer(ctx context.Context, orderID, customerID uuid.UUID) (bool, error)
}

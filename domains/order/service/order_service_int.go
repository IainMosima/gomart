package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/order/schema"
	"github.com/google/uuid"
)

type OrderService interface {
	// Order CRUD operations
	CreateOrder(ctx context.Context, req *schema.CreateOrderRequest) (*schema.CreateOrderResponse, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*schema.OrderResponse, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*schema.OrderResponse, error)
	GetOrderWithDetails(ctx context.Context, orderID uuid.UUID) (*schema.OrderWithDetailsResponse, error)
	UpdateOrder(ctx context.Context, orderID uuid.UUID, req *schema.UpdateOrderRequest) (*schema.OrderResponse, error)
	DeleteOrder(ctx context.Context, orderID uuid.UUID) error

	// Order status management
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, req *schema.UpdateOrderStatusRequest) (*schema.OrderStatusUpdateResponse, error)
	ConfirmOrder(ctx context.Context, orderID uuid.UUID) (*schema.OrderStatusUpdateResponse, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID) (*schema.OrderStatusUpdateResponse, error)
	MarkOrderAsShipped(ctx context.Context, orderID uuid.UUID) (*schema.OrderStatusUpdateResponse, error)
	MarkOrderAsDelivered(ctx context.Context, orderID uuid.UUID) (*schema.OrderStatusUpdateResponse, error)

	// Order listing and search
	ListOrders(ctx context.Context, req *schema.OrderListRequest) (*schema.OrderListResponse, error)
	SearchOrders(ctx context.Context, req *schema.OrderSearchRequest) (*schema.OrderSearchResponse, error)
	GetCustomerOrders(ctx context.Context, customerID uuid.UUID, req *schema.OrderListRequest) (*schema.OrderListResponse, error)
	GetOrdersByStatus(ctx context.Context, status string, req *schema.OrderListRequest) (*schema.OrderListResponse, error)

	// Order items management
	AddOrderItem(ctx context.Context, orderID uuid.UUID, req *schema.AddOrderItemRequest) (*schema.OrderItemResponse, error)
	UpdateOrderItem(ctx context.Context, orderItemID uuid.UUID, req *schema.UpdateOrderItemRequest) (*schema.OrderItemResponse, error)
	RemoveOrderItem(ctx context.Context, orderItemID uuid.UUID) error

	// Order calculations and summary
	GetOrderSummary(ctx context.Context, orderID uuid.UUID) (*schema.OrderSummaryResponse, error)
	RecalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (*schema.OrderResponse, error)

	// Statistics and reporting
	GetOrderStats(ctx context.Context) (*schema.OrderStatsResponse, error)
	GetCustomerOrderStats(ctx context.Context, customerID uuid.UUID) (*schema.OrderStatsResponse, error)

	// Validation
	ValidateOrderRequest(ctx context.Context, req *schema.CreateOrderRequest) error
	CanUpdateOrder(ctx context.Context, orderID uuid.UUID) (bool, error)
	CanCancelOrder(ctx context.Context, orderID uuid.UUID) (bool, error)
}

// NotificationService handles SMS and email notifications (Assignment requirement)
type NotificationService interface {
	// SMS notifications (Africa's Talking requirement)
	SendOrderConfirmationSMS(ctx context.Context, customerPhone, orderNumber string, totalAmount float64) error
	SendOrderStatusUpdateSMS(ctx context.Context, customerPhone, orderNumber, status string) error
	SendOrderDeliveredSMS(ctx context.Context, customerPhone, orderNumber string) error

	// Email notifications to administrators (Assignment requirement)
	SendOrderNotificationEmail(ctx context.Context, order *schema.OrderSummaryResponse) error
	SendOrderStatusUpdateEmail(ctx context.Context, order *schema.OrderSummaryResponse, oldStatus, newStatus string) error
	SendDailyOrderSummaryEmail(ctx context.Context, stats *schema.OrderStatsResponse) error
}

package schema

import (
	"time"

	"github.com/IainMosima/gomart/domains/order/entity"
	"github.com/google/uuid"
)

type OrderResponse struct {
	OrderID         uuid.UUID          `json:"order_id"`
	CustomerID      uuid.UUID          `json:"customer_id"`
	OrderNumber     string             `json:"order_number"`
	Status          entity.OrderStatus `json:"status"`
	TotalAmount     float64            `json:"total_amount"`
	ShippingAddress string             `json:"shipping_address"`
	BillingAddress  string             `json:"billing_address"`
	Notes           *string            `json:"notes,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       *time.Time         `json:"updated_at,omitempty"`
}

type OrderItemResponse struct {
	OrderItemID uuid.UUID `json:"order_item_id"`
	OrderID     uuid.UUID `json:"order_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int32     `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderItemWithDetailsResponse struct {
	OrderItemResponse
	ProductName        string  `json:"product_name"`
	ProductSKU         string  `json:"product_sku"`
	ProductDescription *string `json:"product_description,omitempty"`
}

type OrderWithDetailsResponse struct {
	OrderResponse
	CustomerName  string                          `json:"customer_name"`
	CustomerEmail string                          `json:"customer_email"`
	ItemCount     int                             `json:"item_count"`
	Items         []*OrderItemWithDetailsResponse `json:"items"`
}

type OrderSummaryResponse struct {
	OrderID      uuid.UUID                       `json:"order_id"`
	OrderNumber  string                          `json:"order_number"`
	CustomerID   uuid.UUID                       `json:"customer_id"`
	CustomerName string                          `json:"customer_name"`
	Status       entity.OrderStatus              `json:"status"`
	TotalAmount  float64                         `json:"total_amount"`
	ItemCount    int                             `json:"item_count"`
	Items        []*OrderItemWithDetailsResponse `json:"items"`
	CreatedAt    time.Time                       `json:"created_at"`
	UpdatedAt    *time.Time                      `json:"updated_at,omitempty"`
}

type OrderListResponse struct {
	Orders  []*OrderResponse `json:"orders"`
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
	HasNext bool             `json:"has_next"`
}

type OrderSearchResponse struct {
	Orders  []*OrderWithDetailsResponse `json:"orders"`
	Total   int64                       `json:"total"`
	Page    int                         `json:"page"`
	Limit   int                         `json:"limit"`
	HasNext bool                        `json:"has_next"`
}

type OrderStatsResponse struct {
	TotalOrders       int64   `json:"total_orders"`
	PendingOrders     int64   `json:"pending_orders"`
	CompletedOrders   int64   `json:"completed_orders"`
	CancelledOrders   int64   `json:"cancelled_orders"`
	ProcessingOrders  int64   `json:"processing_orders"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageOrderValue float64 `json:"average_order_value"`
}

type CreateOrderResponse struct {
	Order           *OrderResponse       `json:"order"`
	Items           []*OrderItemResponse `json:"items"`
	PaymentRequired bool                 `json:"payment_required"`
	Message         string               `json:"message"`
}

type OrderStatusUpdateResponse struct {
	OrderID   uuid.UUID          `json:"order_id"`
	OldStatus entity.OrderStatus `json:"old_status"`
	NewStatus entity.OrderStatus `json:"new_status"`
	UpdatedAt time.Time          `json:"updated_at"`
	Message   string             `json:"message"`
}

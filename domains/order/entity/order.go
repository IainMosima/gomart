package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	OrderID         uuid.UUID   `json:"order_id" db:"order_id"`
	CustomerID      uuid.UUID   `json:"customer_id" db:"customer_id"`
	OrderNumber     string      `json:"order_number" db:"order_number"`
	Status          OrderStatus `json:"status" db:"status"`
	TotalAmount     float64     `json:"total_amount" db:"total_amount"`
	ShippingAddress string      `json:"shipping_address" db:"shipping_address"`
	BillingAddress  string      `json:"billing_address" db:"billing_address"`
	Notes           *string     `json:"notes" db:"notes"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time  `json:"updated_at" db:"updated_at"`
	IsDeleted       bool        `json:"is_deleted" db:"is_deleted"`
}

type OrderWithDetails struct {
	Order
	CustomerName  string       `json:"customer_name"`
	CustomerEmail string       `json:"customer_email"`
	ItemCount     int          `json:"item_count"`
	OrderItems    []*OrderItem `json:"order_items,omitempty"`
}

type OrderStats struct {
	TotalOrders       int64   `json:"total_orders"`
	PendingOrders     int64   `json:"pending_orders"`
	CompletedOrders   int64   `json:"completed_orders"`
	CancelledOrders   int64   `json:"cancelled_orders"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageOrderValue float64 `json:"average_order_value"`
}

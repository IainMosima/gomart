package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
)

type Order struct {
	OrderID     uuid.UUID   `json:"order_id" db:"order_id"`
	CustomerID  uuid.UUID   `json:"customer_id" db:"customer_id"`
	OrderNumber string      `json:"order_number" db:"order_number"`
	Status      OrderStatus `json:"status" db:"status"`
	TotalAmount float64     `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at" db:"updated_at"`
	IsDeleted   bool        `json:"is_deleted" db:"is_deleted"`
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	OrderItemID uuid.UUID `json:"order_item_id" db:"order_item_id"`
	OrderID     uuid.UUID `json:"order_id" db:"order_id"`
	ProductID   uuid.UUID `json:"product_id" db:"product_id"`
	Quantity    int32     `json:"quantity" db:"quantity"`
	UnitPrice   float64   `json:"unit_price" db:"unit_price"`
	TotalPrice  float64   `json:"total_price" db:"total_price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
}

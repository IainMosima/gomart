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

type OrderItemWithDetails struct {
	OrderItem
	ProductName        string  `json:"product_name"`
	ProductSKU         string  `json:"product_sku"`
	ProductDescription *string `json:"product_description,omitempty"`
}

type OrderSummary struct {
	OrderID     uuid.UUID               `json:"order_id"`
	OrderNumber string                  `json:"order_number"`
	CustomerID  uuid.UUID               `json:"customer_id"`
	Status      OrderStatus             `json:"status"`
	TotalAmount float64                 `json:"total_amount"`
	ItemCount   int                     `json:"item_count"`
	Items       []*OrderItemWithDetails `json:"items"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   *time.Time              `json:"updated_at,omitempty"`
}

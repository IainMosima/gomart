package schema

import (
	"time"

	"github.com/IainMosima/gomart/domains/order/entity"
	"github.com/google/uuid"
)

type OrderResponse struct {
	OrderID     uuid.UUID          `json:"order_id"`
	CustomerID  uuid.UUID          `json:"customer_id"`
	OrderNumber string             `json:"order_number"`
	Status      entity.OrderStatus `json:"status"`
	TotalAmount float64            `json:"total_amount"`
	CreatedAt   time.Time          `json:"created_at"`
}

type OrderStatusResponse struct {
	OrderID     uuid.UUID          `json:"order_id"`
	OrderNumber string             `json:"order_number"`
	Status      entity.OrderStatus `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
}

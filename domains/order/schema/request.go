package schema

import (
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	CustomerID uuid.UUID                `json:"customer_id" validate:"required"`
	Items      []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int32     `json:"quantity" validate:"required,min=1"`
}

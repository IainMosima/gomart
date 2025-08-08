package schema

import (
	"time"

	"github.com/IainMosima/gomart/domains/order/entity"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	CustomerID      uuid.UUID                `json:"customer_id" validate:"required"`
	ShippingAddress string                   `json:"shipping_address" validate:"required"`
	BillingAddress  string                   `json:"billing_address" validate:"required"`
	Notes           *string                  `json:"notes,omitempty"`
	Items           []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int32     `json:"quantity" validate:"required,min=1"`
}

type UpdateOrderRequest struct {
	Status          *entity.OrderStatus `json:"status,omitempty"`
	ShippingAddress *string             `json:"shipping_address,omitempty"`
	BillingAddress  *string             `json:"billing_address,omitempty"`
	Notes           *string             `json:"notes,omitempty"`
}

type UpdateOrderStatusRequest struct {
	Status entity.OrderStatus `json:"status" validate:"required"`
}

type UpdateOrderItemRequest struct {
	Quantity int32 `json:"quantity" validate:"required,min=1"`
}

type AddOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int32     `json:"quantity" validate:"required,min=1"`
}

type OrderSearchRequest struct {
	CustomerID *uuid.UUID          `json:"customer_id,omitempty"`
	Status     *entity.OrderStatus `json:"status,omitempty"`
	StartDate  *time.Time          `json:"start_date,omitempty"`
	EndDate    *time.Time          `json:"end_date,omitempty"`
	MinAmount  *float64            `json:"min_amount,omitempty" validate:"omitempty,gte=0"`
	MaxAmount  *float64            `json:"max_amount,omitempty" validate:"omitempty,gt=0"`
	Page       int                 `json:"page" validate:"min=1"`
	Limit      int                 `json:"limit" validate:"min=1,max=100"`
}

type OrderListRequest struct {
	CustomerID *uuid.UUID          `json:"customer_id,omitempty"`
	Status     *entity.OrderStatus `json:"status,omitempty"`
	Page       int                 `json:"page" validate:"min=1"`
	Limit      int                 `json:"limit" validate:"min=1,max=100"`
}

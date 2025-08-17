package dtos

import "github.com/google/uuid"

type CreateCategoryRequestDTO struct {
	CategoryName string     `json:"category_name" validate:"required,min=1,max=255"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
}

type UpdateCategoryRequestDTO struct {
	CategoryName string `json:"category_name" validate:"required,min=1,max=255"`
}

type ListCategoriesRequestDTO struct {
	ParentID *uuid.UUID `json:"parent_id,omitempty" form:"parent_id"`
	RootOnly bool       `json:"root_only,omitempty" form:"root_only"`
}

type CreateProductRequestDTO struct {
	ProductName   string    `json:"product_name" validate:"required,min=1,max=255"`
	Description   *string   `json:"description,omitempty" validate:"omitempty,max=1000"`
	Price         float64   `json:"price" validate:"required,gt=0"`
	SKU           string    `json:"sku" validate:"required,min=1,max=100"`
	StockQuantity int32     `json:"stock_quantity" validate:"min=0"`
	CategoryID    uuid.UUID `json:"category_id" validate:"required"`
	IsActive      *bool     `json:"is_active,omitempty"`
}

type UpdateProductRequestDTO struct {
	ProductName   *string    `json:"product_name,omitempty" validate:"omitempty,min=1,max=255"`
	Description   *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	Price         *float64   `json:"price,omitempty" validate:"omitempty,gt=0"`
	StockQuantity *int32     `json:"stock_quantity,omitempty" validate:"omitempty,min=0"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

type UpdateProductStatusRequestDTO struct {
	IsActive bool `json:"is_active"`
}

type UpdateStockRequestDTO struct {
	StockQuantity int32 `json:"stock_quantity" validate:"min=0"`
}

type ProductSearchRequestDTO struct {
	Query      *string    `json:"query,omitempty" form:"query"`
	CategoryID *uuid.UUID `json:"category_id,omitempty" form:"category_id"`
	IsActive   *bool      `json:"is_active,omitempty" form:"is_active"`
	InStock    *bool      `json:"in_stock,omitempty" form:"in_stock"`
	Page       int        `json:"page,omitempty" form:"page" validate:"min=1"`
	Limit      int        `json:"limit,omitempty" form:"limit" validate:"min=1,max=100"`
}

type ValidateTokenRequestDTO struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type CreateOrderRequestDTO struct {
	CustomerID uuid.UUID                   `json:"customer_id" validate:"required"`
	Items      []CreateOrderItemRequestDTO `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItemRequestDTO struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int32     `json:"quantity" validate:"required,min=1"`
}

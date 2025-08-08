package schema

import "github.com/google/uuid"

type CreateProductRequest struct {
	ProductName   string    `json:"product_name" validate:"required,min=1,max=255"`
	Description   *string   `json:"description,omitempty" validate:"omitempty,max=1000"`
	Price         float64   `json:"price" validate:"required,gt=0"`
	SKU           string    `json:"sku" validate:"required,min=1,max=100"`
	StockQuantity int32     `json:"stock_quantity" validate:"min=0"`
	CategoryID    uuid.UUID `json:"category_id" validate:"required"`
	IsActive      *bool     `json:"is_active,omitempty"`
}

type UpdateProductRequest struct {
	ProductName   *string    `json:"product_name,omitempty" validate:"omitempty,min=1,max=255"`
	Description   *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	Price         *float64   `json:"price,omitempty" validate:"omitempty,gt=0"`
	StockQuantity *int32     `json:"stock_quantity,omitempty" validate:"omitempty,min=0"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

type UpdateProductStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type UpdateStockRequest struct {
	StockQuantity int32 `json:"stock_quantity" validate:"min=0"`
}

type BulkUploadProductRequest struct {
	Products []CreateProductRequest `json:"products" validate:"required,min=1,dive"`
}

type ProductSearchRequest struct {
	Query      *string    `json:"query,omitempty"`
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	MinPrice   *float64   `json:"min_price,omitempty" validate:"omitempty,gte=0"`
	MaxPrice   *float64   `json:"max_price,omitempty" validate:"omitempty,gt=0"`
	IsActive   *bool      `json:"is_active,omitempty"`
	InStock    *bool      `json:"in_stock,omitempty"`
	Page       int        `json:"page" validate:"min=1"`
	Limit      int        `json:"limit" validate:"min=1,max=100"`
}

type GetAveragePriceRequest struct {
	CategoryID           uuid.UUID `json:"category_id" validate:"required"`
	IncludeSubcategories bool      `json:"include_subcategories"`
}

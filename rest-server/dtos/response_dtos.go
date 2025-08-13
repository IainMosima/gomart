package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResponseDTO struct {
	CategoryID   uuid.UUID  `json:"category_id"`
	CategoryName string     `json:"category_name"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type CategoryListResponseDTO struct {
	Categories []*CategoryResponseDTO `json:"categories"`
	Total      int64                  `json:"total"`
}

type CategoryAverageProductPriceResponseDTO struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	AveragePrice float64   `json:"average_price"`
}

type ProductResponseDTO struct {
	ProductID     uuid.UUID  `json:"product_id"`
	ProductName   string     `json:"product_name"`
	Description   *string    `json:"description,omitempty"`
	Price         float64    `json:"price"`
	SKU           string     `json:"sku"`
	StockQuantity int32      `json:"stock_quantity"`
	CategoryID    uuid.UUID  `json:"category_id"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type ProductListResponseDTO struct {
	Products []*ProductResponseDTO `json:"products"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	Limit    int                   `json:"limit"`
	HasNext  bool                  `json:"has_next"`
}

type StockUpdateResponseDTO struct {
	ProductID   uuid.UUID `json:"product_id"`
	SKU         string    `json:"sku"`
	OldQuantity int32     `json:"old_quantity"`
	NewQuantity int32     `json:"new_quantity"`
	UpdatedAt   time.Time `json:"updated_at"`
}

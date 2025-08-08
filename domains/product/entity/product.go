package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ProductID     uuid.UUID  `json:"product_id" db:"product_id"`
	ProductName   string     `json:"product_name" db:"product_name"`
	Description   *string    `json:"description" db:"description"`
	Price         float64    `json:"price" db:"price"`
	SKU           string     `json:"sku" db:"sku"`
	StockQuantity int32      `json:"stock_quantity" db:"stock_quantity"`
	CategoryID    uuid.UUID  `json:"category_id" db:"category_id"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted     bool       `json:"is_deleted" db:"is_deleted"`
}

type ProductWithCategory struct {
	Product
	CategoryName string    `json:"category_name" db:"category_name"`
	CategoryPath []*string `json:"category_path,omitempty"`
}

type ProductStats struct {
	TotalProducts   int64   `json:"total_products"`
	ActiveProducts  int64   `json:"active_products"`
	OutOfStockCount int64   `json:"out_of_stock_count"`
	AveragePrice    float64 `json:"average_price"`
	TotalStockValue float64 `json:"total_stock_value"`
}

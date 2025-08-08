package schema

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
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

type ProductWithCategoryResponse struct {
	ProductResponse
	CategoryName string    `json:"category_name"`
	CategoryPath []*string `json:"category_path,omitempty"`
}

type ProductListResponse struct {
	Products []*ProductResponse `json:"products"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	Limit    int                `json:"limit"`
	HasNext  bool               `json:"has_next"`
}

type ProductSearchResponse struct {
	Products []*ProductWithCategoryResponse `json:"products"`
	Total    int64                          `json:"total"`
	Page     int                            `json:"page"`
	Limit    int                            `json:"limit"`
	HasNext  bool                           `json:"has_next"`
	Query    string                         `json:"query,omitempty"`
}

type BulkUploadResponse struct {
	SuccessCount    int                `json:"success_count"`
	FailedCount     int                `json:"failed_count"`
	FailedProducts  []*BulkUploadError `json:"failed_products,omitempty"`
	CreatedProducts []*ProductResponse `json:"created_products,omitempty"`
}

type BulkUploadError struct {
	Index   int                  `json:"index"`
	SKU     string               `json:"sku"`
	Error   string               `json:"error"`
	Product CreateProductRequest `json:"product"`
}

type AveragePriceResponse struct {
	CategoryID            uuid.UUID `json:"category_id"`
	CategoryName          string    `json:"category_name"`
	AveragePrice          float64   `json:"average_price"`
	ProductCount          int64     `json:"product_count"`
	IncludedSubcategories bool      `json:"included_subcategories"`
}

type ProductStatsResponse struct {
	TotalProducts   int64   `json:"total_products"`
	ActiveProducts  int64   `json:"active_products"`
	OutOfStockCount int64   `json:"out_of_stock_count"`
	AveragePrice    float64 `json:"average_price"`
	TotalStockValue float64 `json:"total_stock_value"`
}

type StockUpdateResponse struct {
	ProductID   uuid.UUID `json:"product_id"`
	SKU         string    `json:"sku"`
	OldQuantity int32     `json:"old_quantity"`
	NewQuantity int32     `json:"new_quantity"`
	UpdatedAt   time.Time `json:"updated_at"`
}

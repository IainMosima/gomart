package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/product/schema"
	"github.com/google/uuid"
)

type ProductService interface {
	// Basic CRUD operations
	CreateProduct(ctx context.Context, req *schema.CreateProductRequest) (*schema.ProductResponse, error)
	GetProduct(ctx context.Context, productID uuid.UUID) (*schema.ProductResponse, error)
	GetProductBySKU(ctx context.Context, sku string) (*schema.ProductResponse, error)
	UpdateProduct(ctx context.Context, productID uuid.UUID, req *schema.UpdateProductRequest) (*schema.ProductResponse, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID) error

	// Bulk operations (required by assignment)
	BulkUploadProducts(ctx context.Context, req *schema.BulkUploadProductRequest) (*schema.BulkUploadResponse, error)

	// Product listing and search
	ListProducts(ctx context.Context, req *schema.ProductSearchRequest) (*schema.ProductListResponse, error)
	SearchProducts(ctx context.Context, req *schema.ProductSearchRequest) (*schema.ProductSearchResponse, error)
	GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) (*schema.ProductListResponse, error)

	// Stock management
	UpdateStock(ctx context.Context, productID uuid.UUID, req *schema.UpdateStockRequest) (*schema.StockUpdateResponse, error)
	UpdateProductStatus(ctx context.Context, productID uuid.UUID, req *schema.UpdateProductStatusRequest) (*schema.ProductResponse, error)
	GetLowStockProducts(ctx context.Context, threshold int32) (*schema.ProductListResponse, error)

	// Category-related operations (assignment requirement)
	GetAveragePrice(ctx context.Context, req *schema.GetAveragePriceRequest) (*schema.AveragePriceResponse, error)
	GetAveragePriceByCategory(ctx context.Context, categoryID uuid.UUID) (*schema.AveragePriceResponse, error)

	// Statistics and reporting
	GetProductStats(ctx context.Context) (*schema.ProductStatsResponse, error)
	GetProductStatsByCategory(ctx context.Context, categoryID uuid.UUID) (*schema.ProductStatsResponse, error)

	// Validation
	ValidateProduct(ctx context.Context, req *schema.CreateProductRequest) error
	ValidateSKU(ctx context.Context, sku string) error
}

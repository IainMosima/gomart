package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/product/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	// Basic CRUD operations
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	GetProduct(ctx context.Context, productID uuid.UUID) (*entity.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	UpdateProductStock(ctx context.Context, productID uuid.UUID, stockQuantity int32) (*entity.Product, error)
	UpdateProductStatus(ctx context.Context, productID uuid.UUID, isActive bool) (*entity.Product, error)
	SoftDeleteProduct(ctx context.Context, productID uuid.UUID) error

	// Product listing and filtering
	ListProducts(ctx context.Context) ([]*entity.Product, error)
	ListActiveProducts(ctx context.Context) ([]*entity.Product, error)
	ListProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.Product, error)
	ListProductsInStock(ctx context.Context) ([]*entity.Product, error)

	// Product search and filtering
	SearchProducts(ctx context.Context, query string) ([]*entity.Product, error)
	GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*entity.Product, error)
	GetProductsWithLowStock(ctx context.Context, threshold int32) ([]*entity.Product, error)

	// Category-related operations (for assignment requirement)
	GetAverageProductPrice(ctx context.Context, categoryID uuid.UUID) (float64, error)
	GetAveragePriceByCategoryTree(ctx context.Context, categoryID uuid.UUID) (float64, error)
	GetProductCountByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error)

	// Counting and statistics
	CountProducts(ctx context.Context) (int64, error)
	CountActiveProducts(ctx context.Context) (int64, error)
	CountProductsByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error)

	// Stock management
	DecrementStock(ctx context.Context, productID uuid.UUID, quantity int32) error
	IncrementStock(ctx context.Context, productID uuid.UUID, quantity int32) error

	// Validation helpers
	ProductExists(ctx context.Context, productID uuid.UUID) (bool, error)
	SKUExists(ctx context.Context, sku string) (bool, error)
}

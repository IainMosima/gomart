//go:generate mockgen -source=product_int.go -destination=product_repo_mock.go -package=repository

package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/product/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, product *entity.Product) (*entity.Product, error)
	GetByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) (*entity.Product, error)
	UpdateStock(ctx context.Context, productID uuid.UUID, stockQuantity int32) (*entity.Product, error)
	UpdateStatus(ctx context.Context, productID uuid.UUID, isActive bool) (*entity.Product, error)
	SoftDelete(ctx context.Context, productID uuid.UUID) error

	// Product listing and filtering
	GetAll(ctx context.Context) ([]*entity.Product, error)
	GetActiveProducts(ctx context.Context) ([]*entity.Product, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.Product, error)
	GetInStock(ctx context.Context) ([]*entity.Product, error)

	// Counting and statistics
	CountProducts(ctx context.Context) (int64, error)
	CountActiveProducts(ctx context.Context) (int64, error)
}

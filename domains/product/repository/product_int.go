//go:generate mockgen -source=product_int.go -destination=product_repo_mock.go -package=repository

package repository

import (
	"context"

	"github.com/IainMosima/gomart/domains/product/entity"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) (*entity.Product, error)
	GetByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) (*entity.Product, error)
	Delete(ctx context.Context, productID uuid.UUID) error
	GetAll(ctx context.Context) ([]*entity.Product, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entity.Product, error)
}

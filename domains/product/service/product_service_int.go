//go:generate mockgen -source=product_service_int.go -destination=product_service_mock.go -package=service

package service

import (
	"context"

	"github.com/IainMosima/gomart/domains/product/schema"
	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *schema.CreateProductRequest) (*schema.ProductResponse, error)
	GetProduct(ctx context.Context, productID uuid.UUID) (*schema.ProductResponse, error)
	UpdateProduct(ctx context.Context, productID uuid.UUID, req *schema.UpdateProductRequest) (*schema.ProductResponse, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
	ListProducts(ctx context.Context, req *schema.ProductSearchRequest) (*schema.ProductListResponse, error)
	GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) (*schema.ProductListResponse, error)
}
